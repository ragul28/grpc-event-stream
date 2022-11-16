package psql

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ragul28/grpc-event-stream/pkg/getenv"

	_ "github.com/lib/pq"
)

func CreateConnection() (*sql.DB, error) {

	// postgres db connection vars
	host := getenv.GetEnv("DB_HOST", "localhost")
	user := getenv.GetEnv("DB_USER", "postgres")
	dbname := getenv.GetEnv("DB_NAME", "postgres")
	password := getenv.GetEnv("DB_PASSWORD", "postgres")
	port, _ := strconv.Atoi(getenv.GetEnv("DB_PORT", "5432"))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
