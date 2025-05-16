package structs

type CreateCustomer struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,alphanum,max=255,min=5"`
	LastName string `json:"last_name" validate:"required,alphanum,max=255,min=5"`
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
}

type CheckCustomer struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,max=255,min=5"`
}

type GetCustomer struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Password string `json:"password"`
}

type UpdateCustomer struct {
	Check    CheckCustomer `json:"check" validate:"required"`
	Email    *string       `json:"email,omitempty" validate:"omitempty,email"`
	Name     *string       `json:"name,omitempty" validate:"omitempty,alphanum,max=255,min=5"`
	LastName *string       `json:"last_name,omitempty" validate:"omitempty,alphanum,max=255,min=5"`
}
