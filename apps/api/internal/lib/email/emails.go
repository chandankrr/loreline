package email

func (c *Client) SendWelcomeEmail(to, firstName, baseURL string) error {
	data := map[string]string{
		"UserFirstName": firstName,
		"BaseURL":       baseURL,
	}

	return c.SendEmail(
		to,
		"Welcome to Loreline!",
		TemplateWelcome,
		data,
	)
}

func (c *Client) SendEmailVerification(to, code, baseURL string) error {
	data := map[string]string{
		"VerificationCode": code,
		"ExpiresInMinutes": "15",
		"BaseURL":          baseURL,
	}

	return c.SendEmail(
		to,
		"Confirm your email - Loreline",
		TemplateEmailVerification,
		data,
	)
}

func (c *Client) SendPasswordReset(to, resetURL, baseURL string) error {
	data := map[string]string{
		"ResetURL":         resetURL,
		"ExpiresInMinutes": "30",
		"BaseURL":          baseURL,
	}

	return c.SendEmail(
		to,
		"Reset your password - Loreline",
		TemplatePasswordReset,
		data,
	)
}
