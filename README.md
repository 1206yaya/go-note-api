# Go Note API

This is a RESTful API for managing notes, built with Go, Echo, and DynamoDB.

## Prerequisites

- Go 1.17+
- AWS CLI configured with appropriate credentials
- DynamoDB table created (local or on AWS)

## Setup

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set environment variables:
   - `DYNAMODB_TABLE`: The name of your DynamoDB table
   - `AWS_REGION`: Your AWS region (e.g., "ap-north")
   - `DYNAMODB_ENDPOINT` (optional): For local development, set to "http://localhost:8000"

## Running the API

### Locally

#### Setup DynamoDB

- Uses Docker Compose to run a local DynamoDB instance.
- Tables are created automatically by the application code if they don't exist.

1. For local development, use the provided Docker Compose file to start DynamoDB:

   ```
   make up
   ```

2. Run the API: `make run`

The API will start on `http://localhost:8080`

### On AWS Lambda

1. Build the deployment package: `make package`
2. Upload the `main.zip` file from the `build` directory to your AWS Lambda function
3. Set the Lambda handler to `main`
4. Configure API Gateway to trigger your Lambda function

## Environment Detection

The application automatically detects whether it's running locally or on AWS Lambda:

- When running locally, it loads configuration from a `.env` file (if present)
- On Lambda, it uses the environment variables set in the Lambda configuration

## Initialization

The `initializeApp()` function is used to initialize the application in both environments, ensuring consistency between local and Lambda execution.

## Performance Optimization

When running on Lambda, the application initializes only once and reuses the same instance for subsequent invocations, optimizing performance and reducing cold start times.

## Testing

Run the tests with: `make test`

## API Endpoints

- `POST /notes`: Create a new note
- `GET /notes/:id`: Get a note by ID
- `PUT /notes/:id`: Update a note
- `DELETE /notes/:id`: Delete a note
- `GET /notes`: List all notes

## Security Note

Ensure that you don't log sensitive information like AWS credentials in production environments.

## Error Handling

The application includes error handling for initialization failures, particularly in the Lambda environment. For more detailed error information, check the CloudWatch logs when running on AWS Lambda.
