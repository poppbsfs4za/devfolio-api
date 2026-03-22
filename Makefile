APP=devfolio-api
DB_URL=postgres://postgres:postgres@localhost:5432/devfolio?sslmode=disable

tidy:
	go mod tidy

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

docgen:
	swag init -g cmd/api/main.go --parseDependency --parseInternal

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

docker-build:
	docker build -t $(APP):local .

up:
	docker compose up -d

down:
	docker compose down