package psql

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func CreateConnection() (*sql.DB, error) {

	// postgres db connection vars
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	dbname := getEnv("DB_NAME", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
