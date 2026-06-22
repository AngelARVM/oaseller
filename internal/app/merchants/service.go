package merchants

import (
	"context"
	"fmt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateMerchant(ctx context.Context, req CreateMerchantRequest) (Merchant, error) {
	if req.Name == "" {
		return Merchant{}, fmt.Errorf("name is required")
	}

	merchant := Merchant{
		Name:   req.Name,
		Active: req.Active,
	}

	return s.repository.InsertMerchant(ctx, merchant)
}

func (s *Service) ListMerchants(ctx context.Context) ([]Merchant, error) {
	return s.repository.ListMerchants(ctx)
}

func (s *Service) Merchant(ctx context.Context, merchantId int64) (Merchant, error) {
	return s.repository.Merchant(ctx, merchantId)
}

func (s *Service) PatchMerchant(ctx context.Context, merchantId int64, update UpdateMerchantRequest) (Merchant, error) {
	return s.repository.PatchMerchant(ctx, merchantId, update)
}
