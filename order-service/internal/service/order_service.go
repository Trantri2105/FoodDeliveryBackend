package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"order-service/internal/model"
	"order-service/internal/repository"
	"order-service/pkg/client"
	"order-service/pkg/middleware"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderId int, status string) error
	GetOrderById(ctx context.Context, orderId int) (model.Order, error)
	GetOrderList(ctx context.Context, limit, offset, userId int) ([]model.Order, error)
}

type orderService struct {
	orderRepo        repository.OrderRepository
	deliveryClient   client.DeliveryClient
	restaurantClient client.RestaurantClient
}

const (
	OrderCreated    = "created"
	OrderReady      = "ready for delivery"
	OrderDelivering = "delivering"
	OrderCancelled  = "cancelled"
	DeliveryRange   = 20
)

func (o *orderService) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	order.Status = OrderCreated
	menuItems, err := o.restaurantClient.GetMenu(ctx)
	if err != nil {
		return model.Order{}, err
	}
	order.Subtotal = 0
	for i, orderItem := range order.OrderItems {
		found := false
		for _, menuItem := range menuItems {
			if orderItem.MenuItemId == menuItem.Id {
				if !menuItem.IsAvailable {
					return model.Order{}, errors.New("order item not available")
				}
				found = true
				order.OrderItems[i].UnitPrice = menuItem.Price
				order.OrderItems[i].TotalPrice = menuItem.Price * orderItem.Quantity
				order.Subtotal += order.OrderItems[i].TotalPrice
				break
			}
		}
		if !found {
			return model.Order{}, errors.New("order item not found")
		}
	}
	id, err := o.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return model.Order{}, err
	}
	order.Id = id

	restaurant, err := o.restaurantClient.GetRestaurantInformation(ctx)
	if err != nil {
		return model.Order{}, err
	}

	c := ctx.(*gin.Context)
	token, _ := c.Get(middleware.JWTAccessTokenContextKey)
	accessToken := token.(string)
	delivery, err := o.deliveryClient.CreateDelivery(ctx, order.Id, restaurant.Address, order.ShippingAddress, accessToken)
	if err != nil {
		return model.Order{}, err
	}
	order.Delivery = &model.Delivery{
		Distance:     delivery.Distance,
		Duration:     delivery.Duration,
		Fee:          delivery.Fee,
		FromCoords:   delivery.FromCoords,
		ToCoords:     delivery.ToCoords,
		GeometryLine: delivery.GeometryLine,
		Status:       delivery.Status,
		Shipper: model.Shipper{
			Name:         delivery.Shipper.Name,
			Gender:       delivery.Shipper.Gender,
			Phone:        delivery.Shipper.Phone,
			VehicleType:  delivery.Shipper.VehicleType,
			VehiclePlate: delivery.Shipper.VehiclePlate,
		},
	}
	order.DeliveryFee = delivery.Fee
	order.TotalAmount = order.Subtotal + order.DeliveryFee
	err = o.orderRepo.UpdateOrderFee(ctx, order.Id, order.DeliveryFee, order.TotalAmount)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (o *orderService) UpdateOrderStatus(ctx context.Context, orderId int, status string) error {
	return o.orderRepo.UpdateOrderStatus(ctx, orderId, status)
}

func (o *orderService) GetOrderById(ctx context.Context, orderId int) (model.Order, error) {
	order, err := o.orderRepo.GetOrderById(ctx, orderId)
	if err != nil {
		return model.Order{}, err
	}
	c := ctx.(*gin.Context)
	token, _ := c.Get(middleware.JWTAccessTokenContextKey)
	accessToken := token.(string)
	delivery, err := o.deliveryClient.GetDeliveryByOrderId(ctx, orderId, accessToken)
	if err != nil {
		return model.Order{}, err
	}
	order.Delivery = &model.Delivery{
		Distance:     delivery.Distance,
		Duration:     delivery.Duration,
		Fee:          delivery.Fee,
		FromCoords:   delivery.FromCoords,
		ToCoords:     delivery.ToCoords,
		GeometryLine: delivery.GeometryLine,
		Status:       delivery.Status,
		Shipper: model.Shipper{
			Name:         delivery.Shipper.Name,
			Gender:       delivery.Shipper.Gender,
			Phone:        delivery.Shipper.Phone,
			VehicleType:  delivery.Shipper.VehicleType,
			VehiclePlate: delivery.Shipper.VehiclePlate,
		},
	}
	return order, nil
}

func (o *orderService) GetOrderList(ctx context.Context, limit, offset, userId int) ([]model.Order, error) {
	return o.orderRepo.GetOrderList(ctx, limit, offset, userId)
}

func NewOrderService(orderRepo repository.OrderRepository, deliveryClient client.DeliveryClient, restaurantClient client.RestaurantClient) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		deliveryClient:   deliveryClient,
		restaurantClient: restaurantClient,
	}
}
