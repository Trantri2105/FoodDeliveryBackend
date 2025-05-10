package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type CreateDeliveryRequest struct {
	OrderId           int    `json:"orderId"`
	RestaurantAddress string `json:"restaurantAddress"`
	ShippingAddress   string `json:"shippingAddress"`
}

type DeliveryResponse struct {
	DeliveryId   int             `json:"deliveryId"`
	OrderId      int             `json:"orderId"`
	Distance     float64         `json:"distance"`
	Duration     float64         `json:"duration"`
	Fee          int             `json:"fee"`
	FromCoords   []float64       `json:"fromCoords"`
	ToCoords     []float64       `json:"toCoords"`
	GeometryLine string          `json:"geometryLine"`
	Status       string          `json:"status"`
	Shipper      ShipperResponse `json:"shipper"`
}

type ShipperResponse struct {
	UserId          int    `json:"userId"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Gender          string `json:"gender"`
	Phone           string `json:"phone"`
	Role            string `json:"role"`
	VehicleType     string `json:"vehicleType"`
	VehiclePlate    string `json:"vehiclePlate"`
	TotalDeliveries int    `json:"totalDeliveries"`
	Status          string `json:"status"`
}

type DeliveryClient struct{}

func (d *DeliveryClient) CreateDelivery(ctx context.Context, orderId int, restaurantAddress, shippingAddress, accessToken string) (DeliveryResponse, error) {
	requestURL := os.Getenv("DELIVERY_HOST") + "/delivery"
	log.Printf("delivery url: %s", requestURL)
	r := CreateDeliveryRequest{
		OrderId:           orderId,
		RestaurantAddress: restaurantAddress,
		ShippingAddress:   shippingAddress,
	}
	reqBody, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error marshalling create delivery request body: %v", err)
		return DeliveryResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating create delivery request: %v", err)
		return DeliveryResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending create delivery request: %v", err)
		return DeliveryResponse{}, err
	}
	defer resp.Body.Close()
	log.Println("response status:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return DeliveryResponse{}, fmt.Errorf("create delivery request failed with status code: %d", resp.StatusCode)
	}
	var createResponse DeliveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResponse); err != nil {
		log.Printf("Error parsing create delivery response: %v", err)
		return DeliveryResponse{}, err
	}
	return createResponse, nil
}

func (d *DeliveryClient) GetDeliveryByOrderId(ctx context.Context, orderId int, accessToken string) (DeliveryResponse, error) {
	requestURL := fmt.Sprintf("%s/delivery/order/%d", os.Getenv("DELIVERY_HOST"), orderId)
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		log.Printf("Error creating get delivery by order id request: %v", err)
		return DeliveryResponse{}, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending get delivery by order id request: %v", err)
		return DeliveryResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return DeliveryResponse{}, fmt.Errorf("get delivery by order id request failed with status code: %d", resp.StatusCode)
	}
	var deliveryResponse DeliveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&deliveryResponse); err != nil {
		log.Printf("Error parsing get delivery by order id response: %v", err)
		return DeliveryResponse{}, err
	}
	return deliveryResponse, nil
}
