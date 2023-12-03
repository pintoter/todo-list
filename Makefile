.PHONY: build

build:
	go build -o todo-list ./cmd/app/main.go


.PHONY: run

run: build
	./todo-list

.PHONY: migrations-up

migrations-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

.PHONY: migrations-down

migrations-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down

.PHONY: test

test:
	go test -coverprofile=cover.out -v ./...

test-cover:	go tool cover -html cover.out -o cover.html
	open cover.html

.PHONY: test

swag:
	swag init -g ./cmd/app/main.go