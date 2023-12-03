include .env

.DEFAULT_GOAL = run

MIGRATIONS_DIR = ./migrations
POSTGRES_DSN = postgres://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: build
build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/todo-app ./cmd/app/main.go

.PHONY: run
run: build
	docker-compose up --remove-orphans

rebuild: build
	docker-compose up --remove-orphans --build

.PHONY: stop
stop:
	docker-compose down --remove-orphans

.PHONY: migrations-create
migrations-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) create_notes

.PHONY: migrations-up
migrations-up:
	migrate -path $(MIGRATIONS_DIR) -database $(POSTGRES_DSN) up

.PHONY: migrations-down
migrations-down:
	migrate -path $(MIGRATIONS_DIR) -database $(POSTGRES_DSN) down

.PHONY: test
test:
	go test -coverprofile=cover.out -v ./...
	make --silent test-cover

.PHONY: test-cover
test-cover:	
	go tool cover -html cover.out -o cover.html
	open cover.html

.PHONY: swag
swag:
	swag init -g ./cmd/app/main.go

.PHONY: clean
clean:
	rm -rf ./.bin