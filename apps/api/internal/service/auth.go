package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/chandankrr/loreline/internal/dto"
	"github.com/chandankrr/loreline/internal/lib/job"
	"github.com/chandankrr/loreline/internal/lib/utils/token"
	applogger "github.com/chandankrr/loreline/internal/logger"
	"github.com/chandankrr/loreline/internal/model/user"
	"github.com/chandankrr/loreline/internal/model/verification"
	"github.com/chandankrr/loreline/internal/repository"
	"github.com/chandankrr/loreline/internal/server"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenIssuer       = "loreline"
	accessTokenAudience     = "loreline-api"
	verificationMaxAttempts = 5
	resendCooldown          = 30 * time.Second
)

var (
	ErrInvalidCredentials           = errors.New("invalid email or password")
	ErrInvalidToken                 = errors.New("invalid token")
	ErrExpiredToken                 = errors.New("token has expired")
	ErrEmailInUse                   = errors.New("email already in use")
	ErrUserNotFound                 = errors.New("user not found")
	ErrEmailAlreadyVerified         = errors.New("email is already verified")
	ErrEmailNotVerified             = errors.New("email is not verified")
	ErrInvalidVerificationCode      = errors.New("invalid verification code")
	ErrVerificationExpired          = errors.New("verification code has expired")
	ErrVerificationAttemptsExceeded = errors.New("verification attempts exceeded")
	ErrInvalidResetToken            = errors.New("invalid or expired reset token")
	ErrTooManyRequests              = errors.New("please wait before requesting again")
	ErrOAuthEmailNotVerified        = errors.New("oauth email is not verified")
)

type AccessTokenClaims struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

type AuthService struct {
	server           *server.Server
	userRepo         *repository.UserRepository
	sessionRepo      *repository.SessionRepository
	accountRepo      *repository.AccountRepository
	verificationRepo *repository.VerificationRepository
}

func NewAuthService(
	server *server.Server,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	accountRepo *repository.AccountRepository,
	verificationRepo *repository.VerificationRepository,
) *AuthService {
	return &AuthService{
		server:           server,
		userRepo:         userRepo,
		sessionRepo:      sessionRepo,
		accountRepo:      accountRepo,
		verificationRepo: verificationRepo,
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

	if err := s.sendVerificationEmail(reqCtx, user.Email); err != nil {
		logger.Error().Err(err).Msg("failed to enqueue verification email after registration")
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
) (*user.User, string, string, error) {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	user, err := s.userRepo.GetUserByEmail(reqCtx, payload.Email)
	if err != nil {
		logger.Warn().Msg("authentication failed")
		return nil, "", "", ErrInvalidCredentials
	}

	account, err := s.accountRepo.GetCredentialAccount(reqCtx, user.ID)
	if err != nil {
		// User exists but has no password (might be OAuth only)
		logger.Warn().Msg("authentication failed")
		return nil, "", "", ErrInvalidCredentials
	}

	if err := verifyPassword(*account.Password, payload.Password); err != nil {
		logger.Warn().Msg("authentication failed")
		return nil, "", "", ErrInvalidCredentials
	}

	if !user.EmailVerified {
		logger.Warn().
			Str("user_id", user.ID.String()).
			Msg("login blocked: email is not verified")
		return nil, "", "", ErrEmailNotVerified
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate access token")
		return nil, "", "", err
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
		return nil, "", "", err
	}

	// Business event log
	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "user_login").
		Str("user_id", user.ID.String()).
		Msg("user logged in successfully")

	return user, accessToken, session.Token, nil
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
		return "", "",
			errors.New("email not available from oauth provider; ensure it's public or granted")
	}

	var user *user.User
	isNewUser := false

	// Check if OAuth account already exists
	account, err := s.accountRepo.GetByProviderAndAccountID(reqCtx, gothUser.Provider, gothUser.UserID)
	if err == nil {
		// Existing oauth user
		user, err = s.userRepo.GetUserByID(reqCtx, account.UserID)
		if err != nil {
			logger.Error().
				Err(err).
				Str("provider", gothUser.Provider).
				Msg("failed to get user for existing oauth account")
			return "", "", err
		}

		var expiresAt *time.Time
		if !gothUser.ExpiresAt.IsZero() {
			expiresAt = &gothUser.ExpiresAt
		}

		tx, err := s.server.DB.Pool.Begin(reqCtx)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("failed to begin oauth login transaction")
			return "", "", err
		}
		defer tx.Rollback(reqCtx)

		// Update email verified status if not already verified
		if !user.EmailVerified && isEmailVerifiedByProvider(gothUser) {
			if err := s.userRepo.UpdateEmailVerified(reqCtx, tx, user.ID, true); err != nil {
				logger.Error().
					Err(err).
					Str("provider", gothUser.Provider).
					Msg("failed to update email verified status")
				return "", "", err
			}
		}

		// Optional update tokens
		if _, err := s.accountRepo.UpdateOAuthTokens(
			reqCtx,
			tx,
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

		if err := tx.Commit(reqCtx); err != nil {
			logger.Error().
				Err(err).
				Msg("failed to commit oauth login transaction")
			return "", "", err
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
			if !isEmailVerifiedByProvider(gothUser) {
				logger.Warn().
					Str("provider", gothUser.Provider).
					Str("email", gothUser.Email).
					Msg("oauth email is not verified")
				return "", "", ErrOAuthEmailNotVerified
			}

			tx, err := s.server.DB.Pool.Begin(reqCtx)
			if err != nil {
				logger.Error().Err(err).Msg("failed to begin oauth account linking transaction")
				return "", "", err
			}
			defer tx.Rollback(reqCtx)

			var expiresAt *time.Time
			if !gothUser.ExpiresAt.IsZero() {
				expiresAt = &gothUser.ExpiresAt
			}

			// Link OAuth account
			if _, err := s.accountRepo.CreateOAuthAccount(
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
					Msg("failed to link oauth account to existing user")
				return "", "", err
			}

			if !user.EmailVerified {
				if err := s.userRepo.UpdateEmailVerified(reqCtx, tx, user.ID, true); err != nil {
					logger.Error().
						Err(err).
						Msg("failed to update email verified status during oauth account linking")
				}

				user.EmailVerified = true
			}

			if err := tx.Commit(reqCtx); err != nil {
				logger.Error().Err(err).Msg("failed to commit oauth account linking transaction")
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
) (*user.User, string, string, error) {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	session, err := s.sessionRepo.GetSession(reqCtx, sessionToken)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to get session while refreshing access token")
		return nil, "", "", ErrInvalidToken
	}

	// Check if the session is expired
	if time.Now().UTC().After(session.ExpiresAt) {
		logger.Warn().
			Time("expired_at", session.ExpiresAt).
			Msg("session token has expired")

		if err := s.sessionRepo.RevokeSession(reqCtx, s.server.DB.Pool, sessionToken); err != nil {
			logger.Error().Err(err).Msg("failed to revoke expired session")
		}
		return nil, "", "", ErrExpiredToken
	}

	user, err := s.userRepo.GetUserByID(reqCtx, session.UserID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get user while refreshing access token")
		return nil, "", "", err
	}

	// Generate a new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate access token during refresh")
		return nil, "", "", err
	}

	tx, err := s.server.DB.Pool.Begin(reqCtx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to begin session token rotation transaction")
		return nil, "", "", err
	}
	defer tx.Rollback(reqCtx)

	// Revoke the old session (token rotation)
	if err := s.sessionRepo.RevokeSession(reqCtx, tx, sessionToken); err != nil {
		logger.Error().Err(err).Msg("failed to revoke old session during token rotation")
		return nil, "", "", err
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
		return nil, "", "", err
	}

	if err := tx.Commit(reqCtx); err != nil {
		logger.Error().Err(err).Msg("failed to commit session token rotation")
		return nil, "", "", err
	}

	// Business event log
	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "user_refresh_access_token").
		Str("user_id", user.ID.String()).
		Msg("access token refreshed successfully")

	return user, accessToken, newSession.Token, nil
}

func (s *AuthService) VerifyEmail(ctx echo.Context, email, code string) error {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	user, err := s.userRepo.GetUserByEmail(reqCtx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Str("email", email).Msg("user not found for email verification")
			return ErrInvalidVerificationCode
		}
		logger.Error().
			Err(err).
			Str("email", email).
			Msg("failed to get user by email for verification")
		return err
	}

	if user.EmailVerified {
		logger.Warn().Str("email", email).Msg("email already verified")
		return ErrEmailAlreadyVerified
	}

	record, err := s.verificationRepo.GetLatestByIdentifierAndType(
		reqCtx,
		email,
		verification.TypeEmailVerification,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Str("email", email).Msg("no verification record found")
			return ErrInvalidVerificationCode
		}
		logger.Error().Err(err).Str("email", email).Msg("failed to get verification record")
		return err
	}

	if time.Now().UTC().After(record.ExpiresAt) {
		logger.Warn().Str("email", email).Msg("verification code expired")
		return ErrVerificationExpired
	}

	if record.Attempts >= verificationMaxAttempts {
		logger.Warn().
			Str("email", email).
			Int("attempts", record.Attempts).
			Msg("verification attempts exceeded")
		return ErrVerificationAttemptsExceeded
	}

	codeHash := hashCode(s.server.Config.Auth.VerificationCodeSecret, code)
	if subtle.ConstantTimeCompare([]byte(codeHash), []byte(record.Value)) != 1 {
		attempts, err := s.verificationRepo.IncrementAttempts(
			reqCtx,
			record.ID,
		)
		if err != nil {
			logger.Error().
				Err(err).
				Str("email", email).
				Msg("failed to increment verification attempts")
			return err
		}

		if attempts >= verificationMaxAttempts {
			logger.Warn().
				Str("email", email).
				Int("attempts", attempts).
				Msg("verification attempts exceeded")
			return ErrVerificationAttemptsExceeded
		}

		logger.Warn().Str("email", email).Msg("invalid verification code provided")
		return ErrInvalidVerificationCode
	}

	tx, err := s.server.DB.Pool.Begin(reqCtx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to begin verify email transaction")
		return err
	}
	defer tx.Rollback(reqCtx)

	if err := s.userRepo.UpdateEmailVerified(reqCtx, tx, user.ID, true); err != nil {
		logger.Error().Err(err).Msg("failed to update email_verified")
		return err
	}

	if err := s.verificationRepo.DeleteByIdentifierAndType(
		reqCtx,
		tx,
		email,
		verification.TypeEmailVerification,
	); err != nil {
		logger.Error().Err(err).Msg("failed to delete verification records")
		return err
	}

	if err := tx.Commit(reqCtx); err != nil {
		logger.Error().Err(err).Msg("failed to commit verify email transaction")
		return err
	}

	firstName, _, _ := strings.Cut(strings.TrimSpace(user.Name), " ")

	task, err := job.NewWelcomeEmailTask(user.Email, firstName, s.server.Config.Primary.FrontendURL)
	if err != nil {
		return err
	}

	if s.server.Job != nil && s.server.Job.Client != nil {
		if _, err := s.server.Job.Client.EnqueueContext(reqCtx, task); err != nil {
			return err
		}
	}

	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "email_verified").
		Str("user_id", user.ID.String()).
		Msg("email verified successfully")

	return nil
}

func (s *AuthService) ResendVerificationEmail(ctx echo.Context, email string) error {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	user, err := s.userRepo.GetUserByEmail(reqCtx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Str("email", email).Msg("user not found for resend verification")
			return nil
		}
		logger.Error().
			Err(err).
			Str("email", email).
			Msg("failed to get user by email for resend verification")
		return err
	}

	if user.EmailVerified {
		logger.Warn().Str("email", email).Msg("resend requested for already-verified email")
		return nil
	}

	latest, err := s.verificationRepo.GetLatestByIdentifierAndType(
		reqCtx,
		email,
		verification.TypeEmailVerification,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Error().
			Err(err).
			Str("email", email).
			Msg("failed to get latest verification record")
		return err
	}
	if latest != nil && time.Since(latest.CreatedAt) < resendCooldown {
		logger.Warn().Str("email", email).Msg("verification code requested too quickly")
		return ErrTooManyRequests
	}

	if err := s.sendVerificationEmail(reqCtx, email); err != nil {
		logger.Error().Err(err).Str("email", email).Msg("failed to send verification email")
		return err
	}

	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "verification_email_resent").
		Str("user_id", user.ID.String()).
		Msg("verification email resent")

	return nil
}

func (s *AuthService) ForgotPassword(ctx echo.Context, email string) error {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	user, err := s.userRepo.GetUserByEmail(reqCtx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Str("email", email).Msg("user not found for password reset")
			return nil
		}
		logger.Error().
			Err(err).
			Str("email", email).
			Msg("failed to get user by email for forgot password")
		return err
	}

	// OAuth-only users don't have a credential account yet,
	// but they should still be able to receive a password-reset
	// email and create a password.

	// This flow is therefore for both:
	//	forgot password and set password for OAuth-only account

	latest, err := s.verificationRepo.GetLatestByIdentifierAndType(
		reqCtx,
		email,
		verification.TypePasswordReset,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Error().Err(err).Str("email", email).Msg("failed to get latest reset record")
		return err
	}
	if latest != nil && time.Since(latest.CreatedAt) < resendCooldown {
		logger.Warn().Str("email", email).Msg("password reset requested too quickly")
		return ErrTooManyRequests
	}

	rawToken, err := token.Generate()
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate reset token")
		return err
	}

	hashedToken := token.Hash(rawToken)
	expiresAt := time.Now().UTC().Add(30 * time.Minute)

	if err := s.verificationRepo.DeleteByIdentifierAndType(
		reqCtx,
		s.server.DB.Pool,
		email,
		verification.TypePasswordReset,
	); err != nil {
		logger.Error().
			Err(err).
			Str("email", email).
			Msg("failed to delete old reset records")
		return err
	}

	if _, err := s.verificationRepo.Create(
		reqCtx,
		email,
		hashedToken,
		verification.TypePasswordReset,
		expiresAt,
	); err != nil {
		logger.Error().Err(err).Msg("failed to create reset verification record")
		return err
	}

	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s",
		s.server.Config.Primary.FrontendURL,
		rawToken,
	)
	task, err := job.NewPasswordResetTask(email, resetURL, s.server.Config.Primary.FrontendURL)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create password reset task")
		return err
	}

	if s.server.Job != nil && s.server.Job.Client != nil {
		if _, err := s.server.Job.Client.EnqueueContext(reqCtx, task); err != nil {
			logger.Error().Err(err).Msg("failed to enqueue password reset task")
			return err
		}
	}

	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "forgot_password_requested").
		Str("user_id", user.ID.String()).
		Msg("password reset email enqueued")

	return nil
}

func (s *AuthService) ResetPassword(ctx echo.Context, resetToken, newPassword string) error {
	logger := applogger.GetLogger(ctx)
	reqCtx := ctx.Request().Context()

	hashedToken := token.Hash(resetToken)

	record, err := s.verificationRepo.GetByValueAndType(
		reqCtx,
		hashedToken,
		verification.TypePasswordReset,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Msg("reset token not found")
			return ErrInvalidResetToken
		}
		logger.Error().Err(err).Msg("failed to get verification record for reset token")
		return err
	}

	if time.Now().UTC().After(record.ExpiresAt) {
		logger.Warn().Msg("reset token has expired")
		return ErrInvalidResetToken
	}

	user, err := s.userRepo.GetUserByEmail(reqCtx, record.Identifier)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Str("email", record.Identifier).Msg("user not found for reset token")
			return ErrInvalidResetToken
		}
		logger.Error().Err(err).Msg("failed to get user by email for password reset")
		return err
	}

	newHash, err := hashPassword(newPassword)
	if err != nil {
		logger.Error().Err(err).Msg("failed to hash new password")
		return err
	}

	tx, err := s.server.DB.Pool.Begin(reqCtx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to begin reset password transaction")
		return err
	}
	defer tx.Rollback(reqCtx)

	// Check whether the user already has a credential account
	_, err = s.accountRepo.GetCredentialAccount(reqCtx, user.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Error().
			Err(err).
			Msg("failed to get credential account for password reset")
		return err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// OAuth-only account
		// Create a credential account with the new password
		if _, err := s.accountRepo.CreateCredentialAccount(reqCtx, tx, user.ID, newHash); err != nil {
			logger.Error().
				Err(err).
				Msg("failed to create credential account for password reset")
			return err
		}
	} else {
		// Existing credential account
		// Update the existing password
		if err := s.accountRepo.UpdatePassword(reqCtx, tx, user.ID, newHash); err != nil {
			logger.Error().Err(err).Msg("failed to update account password")
			return err
		}
	}

	if err := s.verificationRepo.DeleteByIdentifierAndType(
		reqCtx,
		tx,
		record.Identifier,
		verification.TypePasswordReset,
	); err != nil {
		logger.Error().Err(err).Msg("failed to delete verification records after password reset")
		return err
	}

	if err := tx.Commit(reqCtx); err != nil {
		logger.Error().Err(err).Msg("failed to commit reset password transaction")
		return err
	}

	// Revoke all existing sessions for this user
	if err := s.sessionRepo.RevokeAllUserSessions(reqCtx, user.ID); err != nil {
		logger.Error().Err(err).Msg("failed to revoke user sessions after password reset")
	}

	eventLogger := applogger.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "password_reset_success").
		Str("user_id", user.ID.String()).
		Msg("password reset successfully")

	return nil
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

func (s *AuthService) sendVerificationEmail(ctx context.Context, email string) error {
	code, err := generateVerificationCode()
	if err != nil {
		return err
	}

	hashedCode := hashCode(s.server.Config.Auth.VerificationCodeSecret, code)
	expiresAt := time.Now().UTC().Add(15 * time.Minute)

	_ = s.verificationRepo.DeleteByIdentifierAndType(
		ctx,
		s.server.DB.Pool,
		email,
		verification.TypeEmailVerification,
	)

	if _, err := s.verificationRepo.Create(
		ctx,
		email,
		hashedCode,
		verification.TypeEmailVerification,
		expiresAt,
	); err != nil {
		return err
	}

	task, err := job.NewEmailVerificationTask(email, code, s.server.Config.Primary.FrontendURL)
	if err != nil {
		return err
	}

	if s.server.Job != nil && s.server.Job.Client != nil {
		if _, err := s.server.Job.Client.EnqueueContext(ctx, task); err != nil {
			return err
		}
	}

	return nil
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

func hashCode(secret, input string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(input))
	return hex.EncodeToString(mac.Sum(nil))
}

func generateVerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
