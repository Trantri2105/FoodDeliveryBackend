package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"restaurant-service/internal/model"
	"restaurant-service/internal/service"
	"restaurant-service/pkg/middleware"
	"strconv"
)

type restaurantHandler struct {
	restaurantService service.RestaurantService
}

func NewRestaurantHandler(restaurantService service.RestaurantService, middleware middleware.AuthMiddleware, r *gin.Engine) {
	h := &restaurantHandler{
		restaurantService: restaurantService,
	}

	restaurant := r.Group("/restaurant")
	restaurant.GET("", h.GetRestaurantInfo)
	restaurant.PATCH("", middleware.ValidateAndExtractJwt(), h.UpdateRestaurantInfo)
	restaurant.POST("/menu/item", middleware.ValidateAndExtractJwt(), h.AddMenuItem)
	restaurant.GET("/menu", h.GetMenu)
	restaurant.PATCH("/menu/item/:id", middleware.ValidateAndExtractJwt(), h.UpdateMenuItem)
	restaurant.DELETE("/menu/item/:id", middleware.ValidateAndExtractJwt(), h.DeleteMenuItem)
	restaurant.GET("/menu/item/:id", h.GetMenuItemById)
}

func (h *restaurantHandler) GetMenuItemById(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id must be integer"})
		return
	}
	res, err := h.restaurantService.GetMenuItemById(c, i)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *restaurantHandler) GetRestaurantInfo(c *gin.Context) {
	restaurantInfo, err := h.restaurantService.GetRestaurantInfo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, restaurantInfo)
}

func (h *restaurantHandler) UpdateRestaurantInfo(c *gin.Context) {
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req model.Restaurant
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.restaurantService.UpdateRestaurantInfo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *restaurantHandler) AddMenuItem(c *gin.Context) {
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req model.MenuItem
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.restaurantService.AddMenuItem(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *restaurantHandler) GetMenu(c *gin.Context) {
	menu, err := h.restaurantService.GetMenu(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, menu)
}

func (h *restaurantHandler) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id must be integer"})
		return
	}
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req model.MenuItem
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = i
	res, err := h.restaurantService.UpdateMenuItem(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *restaurantHandler) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id must be integer"})
		return
	}
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	err = h.restaurantService.DeleteMenuItem(c, i)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "menu item deleted successfully"})
}
