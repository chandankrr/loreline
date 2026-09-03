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
	"github.com/markbates/goth"
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

	user, err := s.userRepo.CreateUser(reqCtx, tx, payload.Name, payload.Email, nil, false)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	_, err = s.accountRepo.CreateCredentialAccount(reqCtx, tx, user.ID, hashPassword)
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

	session, err := s.sessionRepo.CreateSession(reqCtx, s.server.DB.Pool, user.ID, refreshTokenTTL, &ipAddress, &userAgent)
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

func (s *AuthService) OAuthLogin(
	ctx echo.Context,
	gothUser goth.User,
	ipAddress, userAgent string,
) (string, string, error) {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	if gothUser.Email == "" {
		logger.Warn().
			Str("provider", gothUser.Provider).
			Msg("oauth provider returned no email")
		return "", "", errors.New("email not available from oauth provider; ensure it's public or granted")
	}

	var user *user.User
	isNewUser := false

	// Check if OAuth account already exists
	account, err := s.accountRepo.GetByProviderAndAccountID(reqCtx, gothUser.Provider, gothUser.UserID)
	if err == nil {
		// Existing oauth user
		user, err = s.userRepo.GetUserByID(reqCtx, account.UserID)
		if err != nil {
			logger.Error().Err(err).Str("provider", gothUser.Provider).Msg("failed to get user for existing oauth account")
			return "", "", err
		}

		// Optionally update tokens
		var expiresAt *time.Time
		if !gothUser.ExpiresAt.IsZero() {
			expiresAt = &gothUser.ExpiresAt
		}
		if _, err := s.accountRepo.UpdateOAuthTokens(
			reqCtx,
			account.ID,
			&gothUser.AccessToken,
			&gothUser.RefreshToken,
			&gothUser.IDToken,
			expiresAt,
		); err != nil {
			logger.Warn().
				Err(err).
				Str("provider", gothUser.Provider).
				Msg("failed to update oauth tokens")
		}
	} else if errors.Is(err, pgx.ErrNoRows) {
		// Account doesn't exist. Check if user exists by email (Account linking)
		user, err = s.userRepo.GetUserByEmail(reqCtx, gothUser.Email)
		if errors.Is(err, pgx.ErrNoRows) {
			// User doesn't exist either. Create new user + oauth account
			isNewUser = true

			name := gothUser.NickName
			if name == "" {
				name = gothUser.Name
			}

			var expiresAt *time.Time
			if !gothUser.ExpiresAt.IsZero() {
				expiresAt = &gothUser.ExpiresAt
			}

			tx, err := s.server.DB.Pool.Begin(reqCtx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to begin oauth signup transaction")
				return "", "", err
			}
			defer tx.Rollback(reqCtx)

			user, err = s.userRepo.CreateUser(
				reqCtx,
				tx,
				name,
				gothUser.Email,
				&gothUser.AvatarURL,
				isEmailVerifiedByProvider(gothUser))
			if err != nil {
				logger.Error().
					Err(err).
					Str("provider", gothUser.Provider).
					Msg("failed to create user during oauth signup")
				return "", "", err
			}

			if _, err = s.accountRepo.CreateOAuthAccount(
				reqCtx,
				tx,
				user.ID,
				gothUser.Provider,
				gothUser.UserID,
				&gothUser.AccessToken,
				&gothUser.RefreshToken,
				&gothUser.IDToken,
				expiresAt,
				nil,
			); err != nil {
				logger.Error().
					Err(err).
					Str("provider", gothUser.Provider).
					Msg("failed to create oauth account during signup")
				return "", "", err
			}

			if err := tx.Commit(reqCtx); err != nil {
				logger.Error().Err(err).Msg("failed to commit oauth signup transaction")
				return "", "", err
			}
		} else if err != nil {
			logger.Error().Err(err).Msg("failed to get user by email during oauth account linking")
			return "", "", err
		} else {
			// Existing user found by email — link this provider to it
			var expiresAt *time.Time
			if !gothUser.ExpiresAt.IsZero() {
				expiresAt = &gothUser.ExpiresAt
			}

			if _, err := s.accountRepo.CreateOAuthAccount(
				reqCtx,
				s.server.DB.Pool,
				user.ID,
				gothUser.Provider,
				gothUser.UserID,
				&gothUser.AccessToken,
				&gothUser.RefreshToken,
				&gothUser.IDToken,
				expiresAt,
				nil,
			); err != nil {
				logger.Error().
					Err(err).
					Str("provider", gothUser.Provider).
					Msg("failed to link oauth account to existing user")
				return "", "", err
			}
		}
	} else {
		logger.Error().
			Err(err).
			Str("provider", gothUser.Provider).
			Msg("failed to look up oauth account")
		return "", "", err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate access token")
		return "", "", err
	}

	refreshTokenTTL := s.server.Config.Auth.RefreshTokenTTL

	session, err := s.sessionRepo.CreateSession(
		reqCtx,
		s.server.DB.Pool,
		user.ID,
		refreshTokenTTL,
		&ipAddress,
		&userAgent,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create session")
		return "", "", err
	}

	// Business event log
	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "user_oauth_login").
		Str("user_id", user.ID.String()).
		Str("provider", gothUser.Provider).
		Bool("new_user", isNewUser).
		Msg("user logged in via oauth successfully")

	return accessToken, session.Token, nil
}

func (s *AuthService) Logout(ctx echo.Context, sessionToken string) error {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	err := s.sessionRepo.RevokeSession(reqCtx, s.server.DB.Pool, sessionToken)
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

		if err := s.sessionRepo.RevokeSession(reqCtx, s.server.DB.Pool, sessionToken); err != nil {
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
	if err := s.sessionRepo.RevokeSession(reqCtx, tx, sessionToken); err != nil {
		logger.Error().Err(err).Msg("failed to revoke old session during token rotation")
		return "", "", err
	}

	// Issue a new session token
	refreshTokenTTL := s.server.Config.Auth.RefreshTokenTTL

	newSession, err := s.sessionRepo.CreateSession(
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

func isEmailVerifiedByProvider(gothUser goth.User) bool {
	switch gothUser.Provider {
	case "google":
		if v, ok := gothUser.RawData["verified_email"].(bool); ok {
			return v
		}
		return false

	default:
		return false
	}
}
