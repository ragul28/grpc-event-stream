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

docker: build
	docker compose -f ./docker-compose.yml build

docker_update: docker
	docker compose -f ./docker-compose.yml up -d

db_migrate:
	migrate -database ${POSTGRESQL_URL} -path migrations up

gw_local:
	PORT=8083 ORDER_GRPC_ADDR=localhost:8080 go run cmd/gateway/main.go