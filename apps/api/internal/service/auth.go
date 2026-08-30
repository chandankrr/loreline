package service

import (
	"errors"
	"time"

	"github.com/chandankrr/loreline/internal/dto"
	applogger "github.com/chandankrr/loreline/internal/logger"
	"github.com/chandankrr/loreline/internal/model/user"
	"github.com/chandankrr/loreline/internal/repository"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrEmailInUse         = errors.New("email already in use")
)

const (
	accessTokenIssuer   = "loreline"
	accessTokenAudience = "loreline-api"
)

type AccessTokenClaims struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

type AuthService struct {
	server      *server.Server
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	accountRepo *repository.AccountRepository
}

func NewAuthService(
	server *server.Server,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	accountRepo *repository.AccountRepository,
) *AuthService {
	return &AuthService{
		server:      server,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		accountRepo: accountRepo,
	}
}

func (s *AuthService) Register(ctx echo.Context, payload *dto.RegisterPayload) (*user.User, error) {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	_, err := s.userRepo.GetUserByEmail(reqCtx, payload.Email)
	if err == nil {
		logger.Warn().Msg("user already exists with email")
		return nil, ErrEmailInUse
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		logger.Error().Err(err).Msg("failed to get user by email")
		return nil, err
	}

	hashPassword, err := hashPassword(payload.Password)
	if err != nil {
		logger.Error().Err(err).Msg("failed to hash password")
		return nil, err
	}

	tx, err := s.server.DB.Pool.Begin(reqCtx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to begin create user transaction")
		return nil, err
	}
	defer tx.Rollback(reqCtx)

	user, err := s.userRepo.CreateUserTx(reqCtx, tx, payload.Name, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	_, err = s.accountRepo.CreateCredentialAccountTx(reqCtx, tx, user.ID, hashPassword)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create credential account")
		return nil, err
	}

	if err := tx.Commit(reqCtx); err != nil {
		logger.Error().Err(err).Msg("failed to commit create user transaction")
		return nil, err
	}

	// Business event log
	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "user_register").
		Str("user_id", user.ID.String()).
		Msg("user registered successfully")

	return user, nil
}

func (s *AuthService) Login(
	ctx echo.Context,
	payload *dto.LoginPayload,
	ipAddress, userAgent string,
) (string, string, error) {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	user, err := s.userRepo.GetUserByEmail(reqCtx, payload.Email)
	if err != nil {
		logger.Warn().Msg("authentication failed")
		return "", "", ErrInvalidCredentials
	}

	account, err := s.accountRepo.GetCredentialAccount(reqCtx, user.ID)
	if err != nil {
		// User exists but has no password (might be OAuth only)
		logger.Warn().Msg("authentication failed")
		return "", "", ErrInvalidCredentials
	}

	if err := verifyPassword(*account.Password, payload.Password); err != nil {
		logger.Warn().Msg("authentication failed")
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate access token")
		return "", "", err
	}

	refreshTokenTTL := s.server.Config.Auth.RefreshTokenTTL

	session, err := s.sessionRepo.CreateSession(reqCtx, user.ID, refreshTokenTTL, &ipAddress, &userAgent)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create session")
		return "", "", err
	}

	// Business event log
	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "user_login").
		Str("user_id", user.ID.String()).
		Msg("user logged in successfully")

	return accessToken, session.Token, nil
}

func (s *AuthService) Logout(ctx echo.Context, sessionToken string) error {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	err := s.sessionRepo.RevokeSession(reqCtx, sessionToken)
	if err != nil {
		logger.Error().Err(err).Msg("failed to revoke session")
		return err
	}

	// Business event log
	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "user_logout").
		Msg("user logged out")

	return nil
}

func (s *AuthService) RefreshAccessToken(
	ctx echo.Context,
	sessionToken string,
	ipAddress, userAgent string,
) (string, string, error) {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	session, err := s.sessionRepo.GetSession(reqCtx, sessionToken)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to get session while refreshing access token")
		return "", "", ErrInvalidToken
	}

	// Check if the session is expired
	if time.Now().UTC().After(session.ExpiresAt) {
		logger.Warn().Time("expired_at", session.ExpiresAt).Msg("session token has expired")

		if err := s.sessionRepo.RevokeSession(reqCtx, sessionToken); err != nil {
			logger.Error().Err(err).Msg("failed to revoke expired session")
		}
		return "", "", ErrExpiredToken
	}

	user, err := s.userRepo.GetUserByID(reqCtx, session.UserID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user while refreshing access token")
		return "", "", err
	}

	// Generate a new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate access token during refresh")
		return "", "", err
	}

	tx, err := s.server.DB.Pool.Begin(reqCtx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to begin session token rotation transaction")
		return "", "", err
	}
	defer tx.Rollback(reqCtx)

	// Revoke the old session (token rotation)
	if err := s.sessionRepo.RevokeSessionTx(reqCtx, tx, sessionToken); err != nil {
		logger.Error().Err(err).Msg("failed to revoke old session during token rotation")
		return "", "", err
	}

	// Issue a new session token
	refreshTokenTTL := s.server.Config.Auth.RefreshTokenTTL

	newSession, err := s.sessionRepo.CreateSessionTx(
		reqCtx,
		tx,
		user.ID,
		refreshTokenTTL,
		&ipAddress,
		&userAgent,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create new session during token refresh")
		return "", "", err
	}

	if err := tx.Commit(reqCtx); err != nil {
		logger.Error().Err(err).Msg("failed to commit session token rotation")
		return "", "", err
	}

	// Business event log
	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "user_refresh_access_token").
		Str("user_id", user.ID.String()).
		Msg("access token refreshed successfully")

	return accessToken, newSession.Token, nil
}

func (s *AuthService) generateAccessToken(user *user.User) (string, error) {
	now := time.Now().UTC()
	claims := AccessTokenClaims{
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Issuer:    accessTokenIssuer,
			Audience:  jwt.ClaimStrings{accessTokenAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.server.Config.Auth.AccessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.server.Config.Auth.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(_ *jwt.Token) (any, error) {
			return []byte(s.server.Config.Auth.JWTSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(accessTokenIssuer),
		jwt.WithAudience(accessTokenAudience),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func hashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

func verifyPassword(hashedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}
