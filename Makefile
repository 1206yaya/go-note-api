
.PHONY: build test run clean package up down logs

test:
	go test ./...

run:
	go run cmd/main.go

# Frontend commands
run-frontend:
	cd ../react-note-web && npm start

clean:
	rm -rf bin build

# New target for creating the Lambda deployment package
build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/main.go

package: build
	zip build/main.zip bootstrap && rm bootstrap

###	Docker

# Docker Compose commands
up:
	cd docker && docker-compose up -d

down:
	cd docker && docker-compose down

# Cleanup
clean:
	cd docker && docker-compose down -v \
	&& rm -rf dynamodb_data

# Docker Compose logs
logs:
	cd docker && docker-compose logs -f