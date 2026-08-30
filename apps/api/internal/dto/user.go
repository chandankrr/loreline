package dto

import "github.com/go-playground/validator/v10"

type CreateUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=1,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

func (p *CreateUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
