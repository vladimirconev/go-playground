package storage

import (
	"context"

	"example.com/playground/pkg/api"

	"github.com/gofrs/uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type jobOffer struct {
	gorm.Model

	UUID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`

	Company        string      `json:"company"`
	Email          string      `json:"email"`
	ExpirationDate null.String `json:"expiration_date"`
	LinkToOffer    null.String `json:"link"`
	Details        null.String `json:"details"`
	Phone          string      `json:"phone"`
	Salary         float64     `json:"salary"`
}

func (j *jobOffer) TableName() string {
	return "job_offers"
}

type CreateOffer interface {
	Create(context.Context, *api.JobOfferRequest) (*api.JobOfferResponse, error)
}

type GetOffer interface {
	Get(context.Context, string) (*api.JobOfferResponse, error)
}

type GetAllOffers interface {
	GetAll(context.Context, int, int, string) (*api.JobOffersPaginationResponse, error)
}

type UpdateOffer interface {
	Update(context.Context, string, *api.UpdateJobOfferRequest) (*api.JobOfferResponse, error)
}

type DeleteOffer interface {
	DeleteByID(context.Context, string) error
}

type dbService struct {
	db *gorm.DB
}

func (s *dbService) GetAll(ctx context.Context, size, offset int, sortBy string) (*api.JobOffersPaginationResponse, error) {
	var totalCount int64
	var offer jobOffer

	if err := s.db.WithContext(ctx).
		Model(&offer).
		Count(&totalCount).
		Error; err != nil {

		return nil, err
	}

	var offers []jobOffer

	if err := s.db.WithContext(ctx).
		Offset(offset).
		Limit(size).
		Order(sortBy).
		Find(&offers).
		Error; err != nil {

		return nil, err
	}

	data := make([]api.JobOfferResponse, 0)

	for _, o := range offers {
		if o.ID > 0 {
			jobOfferResponse := api.JobOfferResponse{
				ID:             o.UUID.String(),
				Company:        o.Company,
				Email:          o.Email,
				ExpirationDate: o.ExpirationDate.ValueOrZero(),
				LinkToOffer:    o.LinkToOffer.ValueOrZero(),
				Details:        o.Details.ValueOrZero(),
				Salary:         o.Salary,
				ContactPhone:   o.Phone,
			}
			data = append(data, jobOfferResponse)
		}
	}

	return &api.JobOffersPaginationResponse{
		TotalCount: totalCount,
		Data:       data,
	}, nil
}

func (s *dbService) Create(ctx context.Context, req *api.JobOfferRequest) (*api.JobOfferResponse, error) {
	offer := &jobOffer{
		Company:        req.Company,
		Email:          req.Email,
		ExpirationDate: null.StringFrom(req.ExpirationDate),
		LinkToOffer:    null.StringFrom(req.LinkToOffer),
		Details:        null.StringFrom(req.Details),
		Salary:         req.Salary,
		Phone:          req.ContactPhone,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(offer).Error
	})
	if err != nil {
		return nil, err
	}

	return &api.JobOfferResponse{
		ID:             offer.UUID.String(),
		Company:        offer.Company,
		Email:          offer.Email,
		ExpirationDate: offer.ExpirationDate.ValueOrZero(),
		LinkToOffer:    offer.LinkToOffer.ValueOrZero(),
		Details:        offer.Details.ValueOrZero(),
		Salary:         offer.Salary,
		ContactPhone:   offer.Phone,
	}, nil
}

func (s *dbService) Get(ctx context.Context, offerID string) (*api.JobOfferResponse, error) {
	var offer jobOffer

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Where("uuid = (?)", offerID).
			Find(&offer).
			Error
	}); err != nil {
		return nil, err
	}

	return &api.JobOfferResponse{
		ID:             offer.UUID.String(),
		Company:        offer.Company,
		Email:          offer.Email,
		ExpirationDate: offer.ExpirationDate.ValueOrZero(),
		LinkToOffer:    offer.LinkToOffer.ValueOrZero(),
		Details:        offer.Details.ValueOrZero(),
		Salary:         offer.Salary,
		ContactPhone:   offer.Phone,
	}, nil
}

func (s *dbService) Update(ctx context.Context, offerID string, req *api.UpdateJobOfferRequest) (*api.JobOfferResponse, error) {
	var offer jobOffer

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Where("uuid = (?)", offerID).
			Find(&offer).
			UpdateColumn("salary", req.Salary).
			UpdateColumn("email", req.Email).
			UpdateColumn("phone", req.ContactPhone).
			UpdateColumn("link_to_offer", req.LinkToOffer).
			Error
	}); err != nil {
		return nil, err
	}

	return &api.JobOfferResponse{
		ID:             offer.UUID.String(),
		Company:        offer.Company,
		Email:          offer.Email,
		ExpirationDate: offer.ExpirationDate.ValueOrZero(),
		LinkToOffer:    offer.LinkToOffer.ValueOrZero(),
		Details:        offer.Details.ValueOrZero(),
		Salary:         offer.Salary,
		ContactPhone:   offer.Phone,
	}, nil
}

func (s *dbService) DeleteByID(ctx context.Context, offerID string) error {

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Where("uuid = (?)", offerID).
			Delete(&jobOffer{}).
			Error
	})
}

func NewCreateOfferService(db *gorm.DB) CreateOffer { return &dbService{db: db} }

func NewUpdateOfferService(db *gorm.DB) UpdateOffer { return &dbService{db: db} }

func NewGetOfferService(db *gorm.DB) GetOffer { return &dbService{db: db} }

func NewDeleteOfferService(db *gorm.DB) DeleteOffer { return &dbService{db: db} }

func NewGetAllOffersService(db *gorm.DB) GetAllOffers { return &dbService{db: db} }
