package service

import (
	"context"
	"restaurant-service/internal/model"
	"restaurant-service/internal/repository"
)

type RestaurantService interface {
	GetRestaurantInfo(ctx context.Context) (model.Restaurant, error)
	UpdateRestaurantInfo(ctx context.Context, restaurant model.Restaurant) (model.Restaurant, error)
	AddMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error)
	GetMenu(ctx context.Context) ([]model.MenuItem, error)
	UpdateMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error)
	DeleteMenuItem(ctx context.Context, id int) error
}

type restaurantService struct {
	restaurantRepo repository.RestaurantRepository
}

func (r *restaurantService) GetRestaurantInfo(ctx context.Context) (model.Restaurant, error) {
	return r.restaurantRepo.GetRestaurantInfo(ctx)
}

func (r *restaurantService) UpdateRestaurantInfo(ctx context.Context, restaurant model.Restaurant) (model.Restaurant, error) {
	return r.restaurantRepo.UpdateRestaurantInfo(ctx, restaurant)
}

func (r *restaurantService) AddMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error) {
	return r.restaurantRepo.CreateMenuItem(ctx, menuItem)
}

func (r *restaurantService) GetMenu(ctx context.Context) ([]model.MenuItem, error) {
	return r.restaurantRepo.GetMenu(ctx)
}

func (r *restaurantService) UpdateMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error) {
	return r.restaurantRepo.UpdateMenuItem(ctx, menuItem)
}

func (r *restaurantService) DeleteMenuItem(ctx context.Context, id int) error {
	return r.restaurantRepo.DeleteMenuItem(ctx, id)
}

func NewRestaurantService(restaurantRepo repository.RestaurantRepository) RestaurantService {
	return &restaurantService{restaurantRepo: restaurantRepo}
}
