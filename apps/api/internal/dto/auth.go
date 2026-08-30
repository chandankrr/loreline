package dto

import "github.com/go-playground/validator/v10"

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

func (p *RegisterPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

func (p *LoginPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}
