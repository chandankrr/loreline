package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chandankrr/loreline/internal/config"
	"github.com/chandankrr/loreline/internal/lib/email"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

var emailClient *email.Client

func (j *JobService) InitHandlers(config *config.Config, logger *zerolog.Logger) {
	emailClient = email.NewClient(config, logger)
}

func (j *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal welcome email payload: %w", err)
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Processing welcome email task")

	err := emailClient.SendWelcomeEmail(
		p.To,
		p.FirstName,
		p.BaseURL,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "welcome").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Successfully sent welcome email")
	return nil
}

func (j *JobService) handleEmailVerificationTask(ctx context.Context, t *asynq.Task) error {
	var p EmailVerificationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal email verification payload: %w", err)
	}

	j.logger.Info().
		Str("type", "email_verification").
		Str("to", p.To).
		Msg("Processing email verification task")

	err := emailClient.SendEmailVerification(
		p.To,
		p.Code,
		p.BaseURL,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "email_verification").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send email verification")
		return err
	}

	j.logger.Info().
		Str("type", "email_verification").
		Str("to", p.To).
		Msg("Successfully sent email verification")
	return nil
}

func (j *JobService) handlePasswordResetTask(ctx context.Context, t *asynq.Task) error {
	var p PasswordResetPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal password reset payload: %w", err)
	}

	j.logger.Info().
		Str("type", "password_reset").
		Str("to", p.To).
		Msg("Processing password reset task")

	err := emailClient.SendPasswordReset(
		p.To,
		p.ResetURL,
		p.BaseURL,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "password_reset").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send password reset email")
		return err
	}

	j.logger.Info().
		Str("type", "password_reset").
		Str("to", p.To).
		Msg("Successfully sent password reset email")
	return nil
}
