package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UserClient struct{}

// Exported fields for (un)marshalling
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type RegisterUserResponse struct {
	UserID   int64  `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

func (c *UserClient) Register(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error) {
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("USER_SERVICE_URL environment variable is not set")
	}
	requestURL := baseURL + "/auth/register"

	req.Role = "shipper"

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to register user, status code: %d", resp.StatusCode)
	}

	var registerResp RegisterUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &registerResp, nil
}
