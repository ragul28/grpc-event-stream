package repository

import (
	"database/sql"
	"fmt"
	"log"

	pb "github.com/ragul28/grpc-event-stream/event"
)

type Repository interface {
	CreateOrder(*pb.EventRequest) (*pb.EventResponse, error)
	GetOrder(ef *pb.GetEventFilter) (*pb.GetEventResponse, error)
}

type OrderRepository struct {
	DB *sql.DB
}

func (repo *OrderRepository) CreateOrder(er *pb.EventRequest) (*pb.EventResponse, error) {

	sqlStatement := `INSERT INTO orders (id, name) VALUES ($1, $2) RETURNING id`
	id := ""
	err := repo.DB.QueryRow(sqlStatement, er.Id, er.Name).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.EventResponse{Id: er.Id, Success: true}, nil
}

func (repo *OrderRepository) GetOrder(ef *pb.GetEventFilter) (*pb.GetEventResponse, error) {
	var name string

	sqlStatement := `SELECT name from orders where id=$1`
	row := repo.DB.QueryRow(sqlStatement, ef.Id)

	switch err := row.Scan(&name); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("No records no table!")
	case nil:
		log.Println(ef.Id, name)
	default:
		return nil, fmt.Errorf("GetOrder: Bad input :: %e", err)
	}

	return &pb.GetEventResponse{Id: ef.Id, Name: name}, nil
}
