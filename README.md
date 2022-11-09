# gRPC Event stream

Event stream based Go microservice project using Nats JetStream, PostgresSQL & gRPC. Exhibit project with cloud native ecosystem & tools

## Running

* Clone project & fetch go dependency
    ```sh
    go mod download
    ```

* Compile event protocol buffer file & gen gRPC code   
    ```sh
    make proto
    ```

* Run the PSQL & Nats using docker 
    ```sh
    docker compose pull
    docker compose up -d
    ```

* Run all the microservices in seprate shell 
    ```sh
    go run cmd/order-svc/main.go
    go run cmd/payment-svc/main.go
    ```

* Run sample client event flow
    ```sh
    go run cmd/client/main.go
    ```