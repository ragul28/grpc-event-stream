.PHONY: proto build

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    event/orderevent.proto

mod:
	go mod tidy
	go mod verify

build:
	mkdir ./build/bin/ || true
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o=./build/bin/. ./cmd/...

docker:
	docker compose -f ./docker-compose.yml build