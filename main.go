package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Healthy",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}

func setupRoutes(app *fiber.App) {
	app.Get("/", HealthCheck)
}

func initDatabase(hostName string, port int) *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", hostName, port))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	// Database
	const mongodbHostName = "localhost"
	const mongodbHostPort = 27017

	if mongodbClient := initDatabase(mongodbHostName, mongodbHostPort); mongodbClient == nil {
		log.Fatal(fmt.Sprintf("Failed to connect to MongoDB server: %s:%d", mongodbHostName, mongodbHostPort))
	}

	fmt.Println(fmt.Sprintf("Connected to MongoDB server: %s:%d", mongodbHostName, mongodbHostPort))

	// Instantiate Fiber
	app := fiber.New()

	// Fiber Middleware
	app.Use(cors.New())

	// Routes
	setupRoutes(app)

	const appNetworkAddress = ":5000"

	if err := app.Listen(appNetworkAddress); err != nil {
		log.Fatal(err)
	}
}
