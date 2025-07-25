package sqldb

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type SqlInterface interface {
	DB() *sql.DB
}

type sqlStruct struct {
	db *sql.DB
}

func InitPgSql(driver string, host string, port string, username string, password string, database string, schema string) (SqlInterface, error) {
	config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable", host, port, username, password, database, schema)

	DB, err := sql.Open(driver, config)
	if err != nil {
		return nil, err
	}

	DB.SetConnMaxLifetime(time.Minute * 5)
	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(50)

	return &sqlStruct{
		db: DB,
	}, nil

}

func (m *sqlStruct) DB() *sql.DB {
	return m.db
}
