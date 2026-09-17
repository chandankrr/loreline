package email

var PreviewData = map[string]map[string]string{
	"welcome": {
		"UserFirstName": "John",
		"BaseURL":       "https://example.com",
	},

	"password-reset": {
		"BaseURL":          "https://example.com",
		"ResetURL":         "https://example.com/reset-password?token=abc123",
		"ExpiresInMinutes": "30",
	},

	"email-verification": {
		"BaseURL":          "https://example.com",
		"VerificationCode": "123456",
		"ExpiresInMinutes": "15",
	},
}
