package rest

import (
	"net/http"
	"strconv"

	"thanhnt208/delivery-service/api/middlewares"
	"thanhnt208/delivery-service/internal/models"
	"thanhnt208/delivery-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type ShipperHandler struct {
	service services.ShipperService
}

func NewShipperHandler(service services.ShipperService) *ShipperHandler {
	return &ShipperHandler{service: service}
}

// POST /shippers
func (h *ShipperHandler) CreateShipper(c *gin.Context) {
	claims, _ := c.Get(middlewares.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.ShipperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateShipper(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GET /shippers/:id
func (h *ShipperHandler) GetShipperByID(c *gin.Context) {
	claims, _ := c.Get(middlewares.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if !(role == "admin" || role == "shipper") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.service.GetShipperByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GET /shippers?limit=10&offset=0
func (h *ShipperHandler) ListShippers(c *gin.Context) {
	claims, _ := c.Get(middlewares.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	results, err := h.service.ListShippers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
