
.PHONY: build test run clean package



test:
	go test ./...

run:
	go run cmd/main.go

clean:
	rm -rf bin build

# New target for creating the Lambda deployment package
build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/main.go

package: build
	zip build/main.zip bootstrap
# zip -j ../tf-note-infra/lambda/main.zip main