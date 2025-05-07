package structs

type VacancyCreate struct {
	OwnerEmail       string `json:"owner_email" validate:"required,email,max=255"`
	Title            string `json:"title" validate:"required,max=255"`
	DescriptionOffer string `json:"description_offer" validate:"required,min=10"`
	SalaryCents      int    `json:"salary_cents" validate:"required"`
}

type VacancyGet struct {
	Id               uint   `json:"id,omitempty"`
	OwnerEmail       string `json:"owner_email"`
	Title            string `json:"title"`
	DescriptionOffer string `json:"description_offer"`
	SalaryCents      int    `json:"salary_cents"`
	Responses        int    `json:"responses"`
}

type VacancyUpdate struct {
	Title            *string `json:"title" validate:"omitempty,max=255"`
	DescriptionOffer *string `json:"description_offer" validate:"omitempty,min=10"`
	SalaryCents      *int    `json:"salary_cents" validate:"omitempty"`
}
