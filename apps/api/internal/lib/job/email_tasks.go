package job

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TaskWelcome           = "email:welcome"
	TaskEmailVerification = "email:verification"
	TaskPasswordReset     = "email:password_reset"
)

type WelcomeEmailPayload struct {
	To        string `json:"to"`
	FirstName string `json:"first_name"`
	BaseURL   string `json:"base_url"`
}

func NewWelcomeEmailTask(to, firstName, baseURL string) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{
		To:        to,
		FirstName: firstName,
		BaseURL:   baseURL,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskWelcome, payload,
		asynq.MaxRetry(3),
		asynq.Queue("default"),
		asynq.Timeout(30*time.Second)), nil
}

type EmailVerificationPayload struct {
	To      string `json:"to"`
	Code    string `json:"code"`
	BaseURL string `json:"base_url"`
}

func NewEmailVerificationTask(to, code, baseURL string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailVerificationPayload{
		To:      to,
		Code:    code,
		BaseURL: baseURL,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskEmailVerification, payload,
		asynq.MaxRetry(3),
		asynq.Queue("critical"),
		asynq.Timeout(30*time.Second)), nil
}

type PasswordResetPayload struct {
	To       string `json:"to"`
	ResetURL string `json:"reset_url"`
	BaseURL  string `json:"base_url"`
}

func NewPasswordResetTask(to, resetURL, baseURL string) (*asynq.Task, error) {
	payload, err := json.Marshal(PasswordResetPayload{
		To:       to,
		ResetURL: resetURL,
		BaseURL:  baseURL,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskPasswordReset, payload,
		asynq.MaxRetry(3),
		asynq.Queue("critical"),
		asynq.Timeout(30*time.Second)), nil
}
