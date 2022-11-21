package psql

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/ragul28/grpc-event-stream/pkg/utils"
)

func CreateConnection() (*sql.DB, error) {

	// postgres db connection vars
	host := utils.GetEnv("DB_HOST", "localhost")
	user := utils.GetEnv("DB_USER", "postgres")
	dbname := utils.GetEnv("DB_NAME", "postgres")
	password := utils.GetEnv("DB_PASSWORD", "postgres")
	port, _ := strconv.Atoi(utils.GetEnv("DB_PORT", "5432"))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
