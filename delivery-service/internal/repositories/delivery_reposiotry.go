package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"thanhnt208/delivery-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type DeliveryRepository interface {
	CreateDelivery(ctx context.Context, delivery *models.Delivery) (int64, error)
	UpdateDeliveryStatus(ctx context.Context, deliveryID int64, status string) error
	GetDeliveryByID(ctx context.Context, deliveryID int64) (*models.Delivery, error)
	GetDeliveriesByShipperID(ctx context.Context, shipperID int64, limit, offset int) ([]*models.DeliveryGetByShipperId, error)
	GetAvailableShipper(ctx context.Context) (*models.ShipperResponse, error)
	UpdateShipperStatus(ctx context.Context, shipperID int64, status string) error
	GetDeliveryByOrderID(ctx context.Context, orderId int64) (*models.DeliveryResponse, error)
}

type deliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(db *sqlx.DB) DeliveryRepository {
	return &deliveryRepository{db: db}
}

func (r *deliveryRepository) CreateDelivery(ctx context.Context, delivery *models.Delivery) (int64, error) {
	query := `
		INSERT INTO deliveries (
			order_id, shipper_id, restaurant_address, shipping_address, 
			distance, duration, fee, from_coords, to_coords, geometry_line, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING delivery_id
	`

	// Convert coordinates to JSON for storage
	fromCoords, err := json.Marshal(delivery.FromCoords)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal from_coords: %w", err)
	}

	toCoords, err := json.Marshal(delivery.ToCoords)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal to_coords: %w", err)
	}

	var deliveryID int64
	err = r.db.QueryRowContext(
		ctx, query,
		delivery.OrderID,
		delivery.ShipperID,
		delivery.RestaurantAddress,
		delivery.ShippingAddress,
		delivery.Distance,
		delivery.Duration,
		delivery.Fee,
		fromCoords,
		toCoords,
		delivery.GeometryLine,
		delivery.Status,
	).Scan(&deliveryID)

	if err != nil {
		return 0, fmt.Errorf("failed to create delivery: %w", err)
	}

	return deliveryID, nil
}

func (r *deliveryRepository) UpdateDeliveryStatus(ctx context.Context, deliveryID int64, status string) error {
	query := `UPDATE deliveries SET status = $1, updated_at = NOW() WHERE delivery_id = $2`
	result, err := r.db.ExecContext(ctx, query, status, deliveryID)
	if err != nil {
		return fmt.Errorf("failed to update delivery status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no delivery found with id %d", deliveryID)
	}

	if status == "completed" {
		var shipperID int
		err = r.db.GetContext(ctx, &shipperID, "SELECT shipper_id FROM deliveries WHERE delivery_id = $1", deliveryID)
		if err != nil {
			return fmt.Errorf("failed to get shipper_id: %w", err)
		}

		_, err = r.db.ExecContext(ctx, "UPDATE shippers SET total_deliveries = total_deliveries + 1, status = 'available' WHERE id = $1", shipperID)
		if err != nil {
			return fmt.Errorf("failed to update shipper total_deliveries: %w", err)
		}
	}

	return nil
}

func (r *deliveryRepository) GetDeliveryByID(ctx context.Context, deliveryID int64) (*models.Delivery, error) {
	query := `SELECT * FROM deliveries WHERE id = $1`
	var delivery models.Delivery
	err := r.db.GetContext(ctx, &delivery, query, deliveryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("delivery not found")
		}
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	return &delivery, nil
}

func (r *deliveryRepository) GetDeliveriesByShipperID(ctx context.Context, shipperID int64, limit, offset int) ([]*models.DeliveryGetByShipperId, error) {
	query := `
		SELECT delivery_id, order_id, distance, duration, fee, from_coords, to_coords, geometry_line, status
		FROM deliveries
		WHERE shipper_id = $1
		ORDER BY delivery_id DESC
		LIMIT $2 OFFSET $3
	`
	var deliveries []*models.DeliveryGetByShipperId
	err := r.db.SelectContext(ctx, &deliveries, query, shipperID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries by shipper ID: %w", err)
	}
	return deliveries, nil
}

func (r *deliveryRepository) GetAvailableShipper(ctx context.Context) (*models.ShipperResponse, error) {
	query := `
		SELECT id, email, role, name, gender, phone, vehicle_type, vehicle_plate, total_deliveries, status
		FROM shippers
		WHERE status = 'available'
		ORDER BY RANDOM()
		LIMIT 1
	`

	var shipper models.ShipperResponse
	err := r.db.GetContext(ctx, &shipper, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no available shippers found")
		}
		return nil, fmt.Errorf("failed to get available shipper: %w", err)
	}

	return &shipper, nil
}

func (r *deliveryRepository) UpdateShipperStatus(ctx context.Context, shipperID int64, status string) error {
	query := `UPDATE shippers SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, shipperID)
	if err != nil {
		return fmt.Errorf("failed to update shipper status: %w", err)
	}
	return nil
}

func (r *deliveryRepository) GetDeliveryByOrderID(ctx context.Context, orderId int64) (*models.DeliveryResponse, error) {
	query := `
		SELECT * FROM deliveries 
		WHERE order_id = $1
		ORDER BY delivery_id DESC
		LIMIT 1
	`

	var delivery models.Delivery
	err := r.db.GetContext(ctx, &delivery, query, orderId)
	if err != nil {
		return nil, fmt.Errorf("delivery not found: %w", err)
	}

	var shipper models.ShipperResponse
	shipperQuery := `
		SELECT id, email, role, name, gender, phone, vehicle_type, vehicle_plate, total_deliveries, status
		FROM shippers
		WHERE id = $1
	`
	err = r.db.GetContext(ctx, &shipper, shipperQuery, delivery.ShipperID)
	if err != nil {
		return nil, fmt.Errorf("shipper not found: %w", err)
	}

	resp := &models.DeliveryResponse{
		DeliveryID:   delivery.DeliveryID,
		OrderID:      delivery.OrderID,
		Distance:     delivery.Distance,
		Duration:     delivery.Duration,
		Fee:          delivery.Fee,
		FromCoords:   delivery.FromCoords,
		ToCoords:     delivery.ToCoords,
		GeometryLine: delivery.GeometryLine,
		Status:       delivery.Status,
		Shipper:      shipper,
	}

	return resp, nil
}
