package model

type Customer struct {
	Id    uint   `json:"id" validate:"required,email"`
	Email string `json:"email" validate:"required,email"`
}
