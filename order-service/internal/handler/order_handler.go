package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"order-service/internal/model"
	"order-service/internal/service"
	"order-service/pkg/middleware"
	"strconv"
)

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService, middleware middleware.AuthMiddleware, r *gin.Engine) {
	h := &orderHandler{orderService: orderService}
	orderGroup := r.Group("/orders")
	orderGroup.POST("", middleware.ValidateAndExtractJwt(), h.CreateOrder)
	orderGroup.GET("/:id", middleware.ValidateAndExtractJwt(), h.GetOrderById)
	orderGroup.GET("", middleware.ValidateAndExtractJwt(), h.GetOrderList)
	orderGroup.PUT("/:id/status", middleware.ValidateAndExtractJwt(), h.UpdateOrderStatus)
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	userId := int(userClaims["userId"].(float64))
	var req model.Order
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserId = userId
	createdOrder, err := h.orderService.CreateOrder(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdOrder)
}

func (h *orderHandler) GetOrderById(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id must be integer"})
		return
	}
	res, err := h.orderService.GetOrderById(c, i)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *orderHandler) GetOrderList(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	l, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be integer"})
		return
	}
	o, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be integer"})
		return
	}
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	userId := 0
	if role != "admin" {
		userId = int(userClaims["userId"].(float64))
	}
	res, err := h.orderService.GetOrderList(c, l, o, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

type UpdateOrderRequest struct {
	Status string `json:"status"`
}

func (h *orderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id must be integer"})
		return
	}
	var req UpdateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.orderService.UpdateOrderStatus(c, i, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}
