package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1206yaya/go-note-api/internal/handlers"
	"github.com/1206yaya/go-note-api/internal/repositories"
	"github.com/1206yaya/go-note-api/internal/services"
	"github.com/1206yaya/go-note-api/pkg/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func initializeApp() (*echo.Echo, error) {
	// ローカル環境の場合のみ.envを読み込む
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
	}

	log.Println("Environment variables:")
	log.Printf("DYNAMODB_TABLE: %s", os.Getenv("DYNAMODB_TABLE"))
	log.Printf("AWS_REGION: %s", os.Getenv("AWS_REGION"))
	log.Printf("DYNAMODB_ENDPOINT: %s", os.Getenv("DYNAMODB_ENDPOINT"))

	// Initialize DynamoDB client
	dynamoDBClient, err := db.NewDynamoDBClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DynamoDB client: %v", err)
	}

	// 接続テスト
	_, err = dynamoDBClient.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DynamoDB: %v", err)
	}
	log.Println("Successfully connected to DynamoDB")

	err = db.EnsureTableExists(dynamoDBClient, os.Getenv("DYNAMODB_TABLE"))
	if err != nil {
		return nil, fmt.Errorf("failed to ensure table exists: %v", err)
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

	return e, nil
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// アプリケーションの初期化がまだ行われていない場合のみ実行
	if echoLambda == nil {
		// echoLambdaが初期化されていなければ初期化する
		e, err := initializeApp()
		if err != nil {
			log.Printf("Failed to initialize app: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
			}, nil
		}
		echoLambda = echoadapter.New(e)
	}
	// 初期化済みのechoLambdaでリクエストを処理
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running in Lambda
		lambda.Start(Handler)
	} else {
		// Running locally
		e, err := initializeApp()
		if err != nil {
			log.Fatalf("Failed to initialize app: %v", err)
		}
		e.Logger.Fatal(e.Start(":8080"))
	}
}
