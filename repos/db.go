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

func NewStore() (*DB, error) {
	db, err := sql.Open("mysql", "root:mypass@tcp(127.0.0.1:3306)/gts")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
