package services

import (
	"context"
	"fmt"
	"thanhnt208/delivery-service/external/client"
	"thanhnt208/delivery-service/internal/models"
	"thanhnt208/delivery-service/internal/repositories"
)

type ShipperService interface {
	CreateShipper(ctx context.Context, shipper *models.ShipperRequest) (*models.ShipperResponse, error)
	GetShipperByID(ctx context.Context, id int64) (*models.ShipperResponse, error)
	ListShippers(ctx context.Context, limit, offset int) ([]*models.ShipperResponse, error)
}

type shipperService struct {
	repo       repositories.ShipperRepository
	userClient *client.UserClient
}

func NewShipperService(repo repositories.ShipperRepository, userClient *client.UserClient) ShipperService {
	return &shipperService{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *shipperService) CreateShipper(ctx context.Context, shipper *models.ShipperRequest) (*models.ShipperResponse, error) {
	userReq := &client.RegisterUserRequest{
		Email:    shipper.Email,
		Password: shipper.Password,
		Name:     shipper.Name,
		Gender:   shipper.Gender,
		Phone:    shipper.Phone,
		Role:     "shipper",
	}
	userResp, err := s.userClient.Register(ctx, userReq)
	if err != nil {
		return nil, fmt.Errorf("user service register failed: %w", err)
	}

	shipper.ID = userResp.UserID

	return s.repo.Create(ctx, shipper)
}

func (s *shipperService) GetShipperByID(ctx context.Context, id int64) (*models.ShipperResponse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *shipperService) ListShippers(ctx context.Context, limit, offset int) ([]*models.ShipperResponse, error) {
	return s.repo.GetShippers(ctx, limit, offset)
}
