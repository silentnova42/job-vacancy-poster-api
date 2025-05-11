package structs

type ResponseGet struct {
	VacancyId  uint   `json:"vacancy_id,omitempty"`
	Email      string `json:"email"`
	OwnerEmail string `json:"owner_email"`
}

type ResponseCreate struct {
	VacancyId uint   `json:"vacancy_id,omitempty"`
	Email     string `json:"email" validate:"required,email,max=255"`
}
