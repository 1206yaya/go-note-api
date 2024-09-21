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

1. Build the project: `make build`
2. Run the API: `make run`

The API will start on `http://localhost:8080`


## Building for AWS Lambda Deployment

To create a deployment package for AWS Lambda:

```
make deploy-package
```

This will create a `main.zip` file in the `build` directory, which can be used for Lambda deployment.


## Testing

Run the tests with: `make test`

## API Endpoints

- `POST /notes`: Create a new note
- `GET /notes/:id`: Get a note by ID
- `PUT /notes/:id`: Update a note
- `DELETE /notes/:id`: Delete a note
- `GET /notes`: List all notes
