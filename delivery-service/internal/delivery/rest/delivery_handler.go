package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"thanhnt208/delivery-service/api/middlewares"
	"thanhnt208/delivery-service/internal/models"
	"thanhnt208/delivery-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type DeliveryHandler struct {
	service services.DeliveryService
}

func NewDeliveryHandler(service services.DeliveryService) (*DeliveryHandler, error) {
	return &DeliveryHandler{service: service}, nil
}

// POST /delivery
func (h *DeliveryHandler) CreateDelivery(c *gin.Context) {
	claims, _ := c.Get(middlewares.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateDelivery(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DeliveryHandler) UpdateDeliveryStatus(c *gin.Context) {
	claims, _ := c.Get(middlewares.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if !(role == "admin" || role == "shipper") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deliveryIdStr := c.Param("deliveryId")
	deliveryId, err := strconv.ParseInt(deliveryIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deliveryId"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	if err := h.service.UpdateDeliveryStatus(c.Request.Context(), deliveryId, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery status updated successfully"})
}

func (h *DeliveryHandler) GetDeliveriesByShipperID(c *gin.Context) {
	claims, _ := c.Get(middlewares.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if !(role == "admin" || role == "shipper") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	shipperIdStr := c.Param("shipperId")
	shipperId, err := strconv.ParseInt(shipperIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shipperId"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	deliveries, err := h.service.GetDeliveriesByShipperID(c.Request.Context(), shipperId, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deliveries)
}

func (h *DeliveryHandler) GetDeliveryByOrderID(c *gin.Context) {
	claims, _ := c.Get(middlewares.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if !(role == "admin" || role == "customer") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orderIdStr := c.Param("orderId")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid orderId format"})
		return
	}

	deliveries, err := h.service.GetDeliveryByOrderID(c.Request.Context(), orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("delivery for orderId %d not found: %v", orderId, err),
		})
		return
	}

	c.JSON(http.StatusOK, deliveries)
}
