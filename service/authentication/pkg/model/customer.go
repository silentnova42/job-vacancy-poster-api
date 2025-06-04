package model

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
}

type GetCustomer struct {
	Id       uint   `json:"id,omitempty"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}
