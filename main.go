package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var db *sql.DB

	db, err := sql.Open("mysql", "root:mypass@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
}
