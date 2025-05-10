package routes

import (
	"net/http"
	"thanhnt208/delivery-service/api/middlewares"
	"thanhnt208/delivery-service/internal/delivery/rest"
	"thanhnt208/delivery-service/pkg/jwt"
	"thanhnt208/delivery-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	shipperHandler *rest.ShipperHandler,
	deliveryHandler *rest.DeliveryHandler,
	jwtUtils jwt.Utils,
) *gin.Engine {
	router := gin.Default()

	logger := logger.NewLogger("info")

	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.LoggingMiddleware(logger))
	router.Use(middlewares.RecoveryMiddleware(logger))

	authMiddleware := middlewares.NewAuthMiddleware(jwtUtils)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": "1.0.0",
			"status":  "UP",
		})
	})

	shipperGroup := router.Group("/shippers")
	shipperGroup.Use(authMiddleware.ValidateAndExtractJwt())
	setupShipperRoutes(shipperGroup, shipperHandler)

	deliveryGroup := router.Group("/delivery")
	deliveryGroup.Use(authMiddleware.ValidateAndExtractJwt())
	setupDeliveryRoutes(deliveryGroup, deliveryHandler)

	return router
}

func setupShipperRoutes(rg *gin.RouterGroup, handler *rest.ShipperHandler) {
	rg.POST("/", handler.CreateShipper)
	rg.GET("/:id", handler.GetShipperByID)
	rg.GET("/", handler.ListShippers)
}

func setupDeliveryRoutes(rg *gin.RouterGroup, handler *rest.DeliveryHandler) {
	rg.POST("/", handler.CreateDelivery)
	rg.PUT("/:deliveryId/status", handler.UpdateDeliveryStatus)
	rg.GET("/shipper/:shipperId", handler.GetDeliveriesByShipperID)
	rg.GET("/order/:orderId", handler.GetDeliveryByOrderID)
}
