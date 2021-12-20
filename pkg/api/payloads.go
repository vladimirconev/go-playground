package api

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type JobOfferRequest struct {
	Company        string  `json:"company"`
	Email          string  `json:"email"`
	ExpirationDate string  `json:"expiration_date"`
	LinkToOffer    string  `json:"link"`
	Details        string  `json:"details"`
	Salary         float64 `json:"salary"`
	ContactPhone   string  `json:"phone"`
}

func (req JobOfferRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Company, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Details, validation.Required),
		validation.Field(&req.Salary, validation.Required),
		validation.Field(&req.LinkToOffer, validation.Required, is.URL),
		validation.Field(&req.ContactPhone, validation.Required, is.E164),
	)
}

type JobOfferResponse struct {
	ID             string  `json:"uuid"`
	Company        string  `json:"company"`
	Email          string  `json:"email"`
	ExpirationDate string  `json:"expiration_date"`
	LinkToOffer    string  `json:"link"`
	Details        string  `json:"details"`
	Salary         float64 `json:"salary"`
	ContactPhone   string  `json:"phone"`
}

type UpdateJobOfferRequest struct {
	Salary       float64 `json:"salary"`
	Email        string  `json:"email"`
	ContactPhone string  `json:"phone"`
	LinkToOffer  string  `json:"link"`
}

func (req UpdateJobOfferRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Salary, validation.Required),
		validation.Field(&req.LinkToOffer, validation.Required, is.URL),
		validation.Field(&req.ContactPhone, validation.Required, is.E164),
	)
}

type JobOffersPaginationResponse struct {
	TotalCount int64              `json:"total_count"`
	Data       []JobOfferResponse `json:"data"`
}
