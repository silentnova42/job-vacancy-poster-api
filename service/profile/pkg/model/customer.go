package model

type CreateCustomer struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
	Name     string `json:"name" validate:"required,max=255,min=2"`
	LastName string `json:"last_name" validate:"required,max=255,min=2"`
	Resume   string `json:"resume" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
}

type PasswordPayload struct {
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
}

type PasswordUpdateRequest struct {
	OldPassword string `json:"old_password" validate:"required,alphanum,max=255,min=5"`
	NewPassword string `json:"new_password" validate:"required,alphanum,max=255,min=5"`
}

type GetPrivateCustomer struct {
	Id       uint   `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Resume   string `json:"resume"`
}

type GetPublicCustomer struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Resume   string `json:"resume"`
}

type UpdateCustomer struct {
	NewName     *string `json:"new_name,omitempty" validate:"omitempty,alphanum,max=255,min=2"`
	NewLastName *string `json:"new_last_name,omitempty" validate:"omitempty,alphanum,max=255,min=2"`
	NewResume   *string `json:"new_resume,omitempty" validate:"omitempty"`
}
