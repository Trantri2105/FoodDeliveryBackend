package main

import (
	"log"
	"os"

	"thanhnt208/delivery-service/api/routes"
	"thanhnt208/delivery-service/config"
	"thanhnt208/delivery-service/external/client"
	"thanhnt208/delivery-service/internal/delivery/rest"
	"thanhnt208/delivery-service/internal/repositories"
	"thanhnt208/delivery-service/internal/services"
	"thanhnt208/delivery-service/pkg/jwt"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := config.ConnectPostgres()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	mapboxAPIKey := os.Getenv("MAP_BOX_API_KEY")

	shipperRepo := repositories.NewShipperRepository(db)
	userClient := &client.UserClient{}
	shipperService := services.NewShipperService(shipperRepo, userClient)
	shipperHandler := rest.NewShipperHandler(shipperService)

	deliveryRepo := repositories.NewDeliveryRepository(db)
	mapboxClient := client.NewMapboxClient(mapboxAPIKey)
	deliveryService := services.NewDeliveryService(deliveryRepo, mapboxClient)
	deliveryHandler, err := rest.NewDeliveryHandler(deliveryService)

	r := routes.SetupRoutes(shipperHandler, deliveryHandler, jwt.NewJwtUtils())

	port := os.Getenv("DELIVERY_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running at http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
