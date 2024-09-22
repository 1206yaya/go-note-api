# Go Note API

This is a RESTful API for managing notes, built with Go, Echo, and DynamoDB. It is the backend component of the Note App infrastructure. It's designed to work in conjunction with the [react-note-web](https://github.com/1206yaya/react-note-web) frontend and is deployed using the [tf-note-infra](https://github.com/1206yaya/tf-note-infra) Terraform configuration.

## Overview

Backend: Go with Echo framework
Database: Amazon DynamoDB
Deployment: AWS Lambda and API Gateway
Infrastructure as Code: Terraform [tf-note-infra](https://github.com/1206yaya/tf-note-infra)

This API provides a serverless solution for note management, leveraging AWS services for scalability and ease of maintenance. The infrastructure is managed using Terraform, allowing for consistent and version-controlled deployments.

## Table of Contents

- [Go Note API](#go-note-api)
  - [Overview](#overview)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Running the API](#running-the-api)
    - [Locally](#locally)
    - [On AWS Lambda](#on-aws-lambda)
  - [Environment Detection](#environment-detection)
  - [API Endpoints](#api-endpoints)
  - [DynamoDB Admin](#dynamodb-admin)
  - [Testing](#testing)
  - [Security Note](#security-note)
  - [Error Handling](#error-handling)
  - [Performance Optimization](#performance-optimization)

## Prerequisites

- Go 1.17+
- AWS CLI configured with appropriate credentials
- DynamoDB table created (local or on AWS)
- Docker and Docker Compose (for local development)

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-note-api.git
   cd go-note-api
   ```
2. Install dependencies:
   ```
   go mod download
   ```
3. Set up environment variables:
   - Copy the `.env.example` file to `.env`:
     ```
     cp .env.example .env
     ```
   - Edit the `.env` file and set the required environment variables

## Running the API

### Locally

1. Start DynamoDB and DynamoDB Admin:
   ```
   make up
   ```
2. Run the API:
   ```
   make run
   ```

The API will start on `http://localhost:8080`

### On AWS Lambda

1. Build the deployment package:
   ```
   make package
   ```
2. Upload the `main.zip` file from the `build` directory to your AWS Lambda function
3. Set the Lambda handler to `main`
4. Configure API Gateway to trigger your Lambda function

## Environment Detection

The application automatically detects its running environment:

- Local: Loads configuration from a `.env` file (if present)
- AWS Lambda: Uses environment variables set in the Lambda configuration

## API Endpoints

- `POST /notes`: Create a new note
- `GET /notes/:id`: Get a note by ID
- `PUT /notes/:id`: Update a note
- `DELETE /notes/:id`: Delete a note
- `GET /notes`: List all notes

## DynamoDB Admin

For local development, DynamoDB Admin provides a web-based GUI for managing your local DynamoDB instance.

- Access: `http://localhost:8001`
- Features:
  - View and edit tables
  - Run queries
  - Manage indexes
  - Import/export data

## Testing

Run the tests with:

```
make test
```

## Security Note

Ensure that you don't log sensitive information like AWS credentials in production environments.

## Error Handling

The application includes comprehensive error handling. For detailed error information when running on AWS Lambda, check the CloudWatch logs.

## Performance Optimization

When running on Lambda, the application initializes once and reuses the same instance for subsequent invocations, optimizing performance and reducing cold start times.
