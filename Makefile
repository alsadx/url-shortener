.PHONY:
.SILENT:

url:
	go build -o url-shortener ./cmd/app/main.go
	./url-shortener

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app

test:
	go test -v ./...

swag:
	swag init -g internal/app/app.go

lint:
	golangci-lint run