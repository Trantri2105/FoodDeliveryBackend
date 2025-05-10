package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"order-service/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) (int, error)
	UpdateOrderStatus(ctx context.Context, orderId int, status string) error
	GetOrderById(ctx context.Context, orderId int) (model.Order, error)
	GetOrderList(ctx context.Context, limit, offset, userId int) ([]model.Order, error)
	UpdateOrderFee(ctx context.Context, orderId, deliveryFee, total int) error
}

type orderRepository struct {
	db *sqlx.DB
}

func (o *orderRepository) UpdateOrderFee(ctx context.Context, orderId, deliveryFee, total int) error {
	query := `UPDATE orders SET delivery_fee=$1, total_amount=$2 WHERE id=$3`
	_, err := o.db.ExecContext(ctx, query, deliveryFee, total, orderId)
	if err != nil {
		log.Printf("Order repo, error update order fee: %v", err)
		return err
	}
	return nil
}

func (o *orderRepository) CreateOrder(ctx context.Context, order model.Order) (orderId int, err error) {
	var tx *sqlx.Tx
	tx, err = o.db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO orders (user_id, shipping_address, phone_number, status, subtotal, delivery_fee, total_amount)
				VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = tx.QueryRowContext(ctx, query, order.UserId, order.ShippingAddress, order.PhoneNumber, order.Status, order.Subtotal, order.DeliveryFee, order.TotalAmount).Scan(&orderId)
	if err != nil {
		log.Printf("Order repo, error creating order and getting ID: %v", err)
		return 0, err // Return zero ID and the error
	}
	err = tx.GetContext(ctx, &orderId, query, order)
	for i := range order.OrderItems {
		order.OrderItems[i].OrderId = orderId
	}

	query = `INSERT INTO order_items (menu_item_id, quantity, unit_price, total_price, order_id) 
			VALUES (:menu_item_id, :quantity, :unit_price, :total_price, :order_id)`

	_, err = tx.NamedExecContext(ctx, query, order.OrderItems)
	if err != nil {
		log.Printf("Order repo, error creating order items: %v", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Order repo, error committing transaction: %v", err)
		return
	}

	return
}

func (o *orderRepository) UpdateOrderStatus(ctx context.Context, orderId int, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	res, err := o.db.ExecContext(ctx, query, status, orderId)
	if err != nil {
		log.Printf("Order repo, error update status: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Order repo, error getting rows affected: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (o *orderRepository) GetOrderById(ctx context.Context, orderId int) (model.Order, error) {
	query := `SELECT * FROM orders WHERE id = $1`
	row := o.db.QueryRowxContext(ctx, query, orderId)
	var order model.Order
	err := row.StructScan(&order)
	if err != nil {
		log.Printf("Order repo, error get order: %v", err)
		return order, err
	}
	query = `SELECT menu_item_id, quantity, unit_price, total_price FROM order_items WHERE order_id = $1`
	rows, err := o.db.QueryxContext(ctx, query, orderId)
	if err != nil {
		log.Printf("Order repo, error get order_items: %v", err)
		return order, err
	}
	for rows.Next() {
		var orderItem model.OrderItem
		err = rows.StructScan(&orderItem)
		if err != nil {
			log.Printf("Order repo, error get order_items: %v", err)
			return order, err
		}
		order.OrderItems = append(order.OrderItems, orderItem)
	}
	return order, err

}

func (o *orderRepository) GetOrderList(ctx context.Context, limit, offset, userId int) ([]model.Order, error) {
	var query string
	if userId != 0 {
		query = fmt.Sprintf("SELECT * FROM orders WHERE user_id = %d ORDER BY id DESC LIMIT $1 OFFSET $2", userId)
	} else {
		query = `SELECT * FROM orders ORDER BY id DESC LIMIT $1 OFFSET $2`
	}
	rows, err := o.db.QueryxContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("Order repo, error get order list: %v", err)
		return nil, err
	}
	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err = rows.StructScan(&order)
		if err != nil {
			log.Printf("Order repo, error get order list: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, err

}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}
