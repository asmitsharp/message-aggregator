package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ashmitsharp/msg-agg/internal/handlers"
	"github.com/ashmitsharp/msg-agg/internal/models"
	"github.com/ashmitsharp/msg-agg/internal/routes"
	"github.com/ashmitsharp/msg-agg/internal/services"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectWithRetry(connectFunc func() error, service string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		err := connectFunc()
		if err == nil {
			log.Printf("Connected to %s", service)
			return nil
		}
		log.Printf("Failed to connect to %s (attempt %d/%d): %v", service, i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("failed to connect to %s after %d attempts", service, maxRetries)
}

func main() {
	maxRetries := 5

	// Database connection
	var db *gorm.DB
	err := connectWithRetry(func() error {
		var err error
		dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		return err
	}, "PostgreSQL", maxRetries)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.MasterUser{}, &models.UserIdentity{}, &models.User{})

	// Redis connection
	var rdb *redis.Client
	err = connectWithRetry(func() error {
		rdb = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		})
		return rdb.Ping(context.Background()).Err()
	}, "Redis", maxRetries)
	if err != nil {
		log.Fatal(err)
	}

	// RabbitMQ connection
	var conn *amqp.Connection
	err = connectWithRetry(func() error {
		var err error
		conn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
		return err
	}, "RabbitMQ", maxRetries)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Synapse connection
	synapseURL := fmt.Sprintf("%s/_matrix/client/versions", os.Getenv("MATRIX_HOMESERVER"))
	err = connectWithRetry(func() error {
		_, err := http.Get(synapseURL)
		return err
	}, "Synapse", maxRetries)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	platformAuthService := services.NewPlatformAuthService()
	userService := services.NewUserService(db, platformAuthService)
	authHandler := handlers.NewAuthHandler(userService)
	identityHandler := handlers.NewIdentityHandler(userService)

	routes.SetupRoutes(app, authHandler, identityHandler)

	app.Listen(":3000")
}
