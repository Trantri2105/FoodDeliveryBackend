package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type RestaurantInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	IsActive    bool   `json:"isActive"`
	OpenTime    string `json:"openTime"`
	CloseTime   string `json:"closeTime"`
}

type MenuItemResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	IsAvailable bool   `json:"isAvailable"`
}

type RestaurantClient struct {
}

func (c *RestaurantClient) GetRestaurantInformation(ctx context.Context) (RestaurantInfoResponse, error) {
	requestURL := os.Getenv("RESTAURANT_HOST") + "/restaurant"
	// Create a request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return RestaurantInfoResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending get restaurant info request: %v", err)
		return RestaurantInfoResponse{}, err
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading get restaurant info response: %v", err)
		return RestaurantInfoResponse{}, err
	}

	var restaurantInfo RestaurantInfoResponse
	err = json.Unmarshal(responseData, &restaurantInfo)
	if err != nil {
		log.Printf("Error parsing get restaurant info response: %v", err)
		return RestaurantInfoResponse{}, err
	}

	return restaurantInfo, nil
}

func (c *RestaurantClient) GetMenu(ctx context.Context) ([]MenuItemResponse, error) {
	requestURL := os.Getenv("RESTAURANT_HOST") + "/restaurant/menu"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending get menu request: %v", err)
		return nil, err
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading get menu response: %v", err)
		return nil, err
	}

	var menuItems []MenuItemResponse
	err = json.Unmarshal(responseData, &menuItems)
	if err != nil {
		log.Printf("Error parsing get menu response: %v", err)
		return nil, err
	}
	return menuItems, nil
}
