.PHONY: build

build:
	go build -o todo-list ./cmd/app/main.go


.PHONY: run

run: build
	./todo-list