proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    event/orderevent.proto

mod:
	go mod tidy
	go mod verify