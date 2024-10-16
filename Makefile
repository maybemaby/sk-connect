PATH := $(PATH):$(PWD)/app/node_modules/.bin

dev:
	APP_ENV=development ALLOWED_HOSTS=* air

test:
	go test -v ./...

lint:
	buf lint

generate:
	buf generate

build:
	go build -o bin/api cmd/api/main.go
