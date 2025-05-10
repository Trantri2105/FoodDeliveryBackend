package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"restaurant-service/internal/handler"
	"restaurant-service/internal/repository"
	"restaurant-service/internal/service"
	"restaurant-service/pkg/jwt"
	"restaurant-service/pkg/middleware"
	"time"

	_ "github.com/lib/pq"
)

func loadEnvVariable() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func PostgresConnect() *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("DB_NAME"))
	log.Print("Connecting to PostgreSQL: ", psqlInfo)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db
}

func main() {
	loadEnvVariable()
	db := PostgresConnect()
	restaurantRepo := repository.NewRestaurantRepository(db)
	utils := jwt.NewJwtUtils()
	restaurantService := service.NewRestaurantService(restaurantRepo)
	m := middleware.NewAuthMiddleware(utils)
	r := gin.Default()
	handler.NewRestaurantHandler(restaurantService, m, r)
	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
