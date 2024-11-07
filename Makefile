APP_NAME?=connect-authz

.PHONY: all
all: build test clean

.PHONY: build
build: 
	@go build -o ./bin/ ./...

.PHONY: test
test: 
	@go test -v -count=1 -cover ./...

test/cover:
	@mkdir -p "./bin"
	@go test -short -coverprofile=bin/cov.out `go list ./... | grep -v vendor/`
	@go tool cover -func=bin/cov.out

.PHONY: clean
clean:
	@rm -rf ./bin

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: validate
validate: 
	@go run ./cmd/schema-validator/...

.PHONY: run
run:
	@go run ./cmd/authz/...

.PHONY: docker-build
docker-build:
	@docker build --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t $(APP_NAME) .

.PHONY: up
up:
	@docker compose -f ./docker/docker-compose.yml -p authz up -d --build --remove-orphans

.PHONY: down
down:
	@docker compose -p authz down                                   

deps:
	@go install github.com/matryer/moq@v0.4.0
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

.PHONY: generate
generate:
	@go generate ./...

## Build API documentation 
.PHONY: docs
docs: 
	swag fmt -d ./
	swag init -g internal/app/authz/authz.go
