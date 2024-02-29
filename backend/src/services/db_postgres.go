package services

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" 
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewDBConnection(config DBConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func QueryDB(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	return rows, nil
}

func ExecDB(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	return result, nil
}

func QueryRowDB(db *sql.DB, query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
