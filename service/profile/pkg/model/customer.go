package model

type CreateCustomer struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,max=255,min=2"`
	LastName string `json:"last_name" validate:"required,max=255,min=2"`
	Resume   string `json:"resume" validate:"required"`
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
}

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
}

type GetPrivateCustomer struct {
	Id       uint   `json:"id,omitempty"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Resume   string `json:"resume"`
	Password string `json:"password"`
}

type GetPublicCustomer struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Resume   string `json:"resume"`
}

type UpdateCustomer struct {
	Credentials Credentials `json:"credentials" validate:"required"`
	Email       *string     `json:"email,omitempty" validate:"omitempty,email"`
	Name        *string     `json:"name,omitempty" validate:"omitempty,alphanum,max=255,min=2"`
	LastName    *string     `json:"last_name,omitempty" validate:"omitempty,alphanum,max=255,min=2"`
	Resume      *string     `json:"resume,omitempty" validate:"omitempty"`
	Password    *string     `json:"password" validate:"required,alphanum,max=255,min=5"`
}
