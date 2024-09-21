package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/1206yaya/go-note-api/internal/handlers"
	"github.com/1206yaya/go-note-api/internal/repositories"
	"github.com/1206yaya/go-note-api/internal/services"
	"github.com/1206yaya/go-note-api/pkg/db"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	log.Println("Environment variables:")
	log.Printf("DYNAMODB_TABLE: %s", os.Getenv("DYNAMODB_TABLE"))
	log.Printf("AWS_REGION: %s", os.Getenv("AWS_REGION"))
	log.Printf("AWS_ACCESS_KEY_ID: %s", os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Printf("AWS_SECRET_ACCESS_KEY: %s", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	log.Printf("DYNAMODB_ENDPOINT: %s", os.Getenv("DYNAMODB_ENDPOINT"))
}

func main() {

	// Initialize DynamoDB client
	dynamoDBClient, err := db.NewDynamoDBClient()
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}
	// 接続テスト
	_, err = dynamoDBClient.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatalf("Failed to connect to DynamoDB: %v", err)
	}

	log.Println("Successfully connected to DynamoDB")

	err = db.EnsureTableExists(dynamoDBClient, os.Getenv("DYNAMODB_TABLE"))
	if err != nil {
		log.Fatalf("Failed to ensure table exists: %v", err)
	}

	// Initialize repository
	noteRepo := repositories.NewDynamoDBNoteRepository(dynamoDBClient, os.Getenv("DYNAMODB_TABLE"))

	// Initialize service
	noteService := services.NewNoteService(noteRepo)

	// Initialize handler
	noteHandler := handlers.NewNoteHandler(noteService)

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// FE_URL: Frontendの本番環境のドメイン
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/notes", noteHandler.CreateNote)
	e.GET("/notes/:id", noteHandler.GetNote)
	e.PUT("/notes/:id", noteHandler.UpdateNote)
	e.DELETE("/notes/:id", noteHandler.DeleteNote)
	e.GET("/notes", noteHandler.ListNotes)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
