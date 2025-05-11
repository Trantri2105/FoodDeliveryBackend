package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
	"restaurant-service/internal/model"
	"strings"
)

type RestaurantRepository interface {
	GetRestaurantInfo(ctx context.Context) (model.Restaurant, error)
	UpdateRestaurantInfo(ctx context.Context, restaurant model.Restaurant) (model.Restaurant, error)
	CreateMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error)
	GetMenu(ctx context.Context) ([]model.MenuItem, error)
	UpdateMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error)
	DeleteMenuItem(ctx context.Context, id int) error
}

type restaurantRepository struct {
	db *sqlx.DB
}

func (r *restaurantRepository) GetRestaurantInfo(ctx context.Context) (model.Restaurant, error) {
	query := `SELECT name, description, address, phone_number, is_active, open_time, close_time FROM restaurants`

	row := r.db.QueryRowxContext(ctx, query)
	var restaurant model.Restaurant
	err := row.StructScan(&restaurant)
	if err != nil {
		log.Printf("Restaurant repo, get restaurant info err: %v", err)
	}
	return restaurant, err
}

func (r *restaurantRepository) UpdateRestaurantInfo(ctx context.Context, restaurant model.Restaurant) (model.Restaurant, error) {
	var updateFields []string
	var values []interface{}
	t := reflect.TypeOf(restaurant)
	v := reflect.ValueOf(restaurant)
	cnt := 1
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if !value.IsZero() {
			updateFields = append(updateFields, fmt.Sprintf("%s = $%d", field.Tag.Get("db"), cnt))
			values = append(values, value.Interface())
			cnt += 1
		}
	}
	if len(updateFields) == 0 {
		return model.Restaurant{}, nil
	}
	query := fmt.Sprintf("UPDATE restaurants SET %s WHERE id = 1 RETURNING  name, description, address, phone_number, is_active, open_time, close_time", strings.Join(updateFields, ", "))
	row := r.db.QueryRowxContext(ctx, query, values...)
	var updatedRestaurant model.Restaurant
	err := row.StructScan(&updatedRestaurant)
	if err != nil {
		log.Printf("Restaurant repository, update restaurant err: %v", err)
		return model.Restaurant{}, err
	}
	return updatedRestaurant, nil
}

func (r *restaurantRepository) CreateMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error) {
	query := `INSERT INTO menu_items (restaurant_id, name, description, price, is_available, image_url)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, description, price, is_available, image_url`
	row := r.db.QueryRowxContext(ctx, query, 1, menuItem.Name, menuItem.Description, menuItem.Price, menuItem.IsAvailable, menuItem.ImageUrl)
	var newMenuItem model.MenuItem
	err := row.StructScan(&newMenuItem)
	if err != nil {
		log.Printf("Restaurant repo, create menu item err: %v", err)
		return model.MenuItem{}, err
	}
	return newMenuItem, nil
}

func (r *restaurantRepository) GetMenu(ctx context.Context) ([]model.MenuItem, error) {
	query := `SELECT id, name, description, price, is_available, image_url FROM menu_items ORDER BY id`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		log.Printf("Restaurant repo, error getting menu items: %v", err)
		return nil, err
	}
	var menu []model.MenuItem
	for rows.Next() {
		var item model.MenuItem
		err = rows.StructScan(&item)
		if err != nil {
			log.Printf("Restaurant repo, error getting menu items: %v", err)
			return nil, err
		}
		menu = append(menu, item)
	}
	return menu, nil
}

func (r *restaurantRepository) UpdateMenuItem(ctx context.Context, menuItem model.MenuItem) (model.MenuItem, error) {
	var updateFields []string
	var values []interface{}
	t := reflect.TypeOf(menuItem)
	v := reflect.ValueOf(menuItem)
	cnt := 1
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if !value.IsZero() {
			updateFields = append(updateFields, fmt.Sprintf("%s = $%d", field.Tag.Get("db"), cnt))
			values = append(values, value.Interface())
			cnt += 1
		}
	}
	if len(updateFields) == 0 {
		return model.MenuItem{}, nil
	}
	values = append(values, menuItem.Id)
	query := fmt.Sprintf("UPDATE menu_items SET %s WHERE id = $%d RETURNING id, name, description, price, is_available, image_url", strings.Join(updateFields, ", "), cnt)
	row := r.db.QueryRowxContext(ctx, query, values...)
	var updatedMenuItem model.MenuItem
	err := row.StructScan(&updatedMenuItem)
	if err != nil {
		log.Printf("Restaurant repository, update menu item err: %v", err)
		return model.MenuItem{}, err
	}
	return updatedMenuItem, nil
}

func (r *restaurantRepository) DeleteMenuItem(ctx context.Context, id int) error {
	query := `DELETE FROM menu_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Restaurant repo, delete menu item err: %v", err)
	}
	return err
}

func NewRestaurantRepository(db *sqlx.DB) RestaurantRepository {
	return &restaurantRepository{db: db}
}
