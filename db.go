package main

import (
	"database/sql"
	"log"
	"os"
)

func create_database() {
	var err error
	db, err = sql.Open("mysql", "root:"+os.Getenv("DB_PASSWORD")+"@localhost/proj")
	if err != nil {
		log.Fatal(err)
	}
}
