package model

type Customer struct {
	Email string `json:"email" validate:"required,email"`
}
