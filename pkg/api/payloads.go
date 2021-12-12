package api

type JobOfferRequest struct {
	Company        string  `json:"company"`
	Email          string  `json:"email"`
	ExpirationDate string  `json:"expiration_date"`
	LinkToOffer    string  `json:"link"`
	Details        string  `json:"details"`
	Salary         float64 `json:"salary"`
	ContactPhone   string  `json:"phone"`
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
