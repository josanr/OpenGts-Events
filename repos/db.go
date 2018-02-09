package repos

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

type Rows struct {
	*sql.Rows
}

func NewStore(login, pass, database, host, port string) (*DB, error) {
	db, err := sql.Open("mysql", login + ":" + pass + "@tcp(" + host + ":" + port + ")/" + database)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
