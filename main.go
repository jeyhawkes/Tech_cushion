package main

import (
	"net/http"

	"github.com/jeyhawkes/tech_cushion/database"
)

const (
	db_username = "root"
	db_password = "password"
	db_name     = "cushion"
)

func main() {

	var db database.Database
	/*
		if err := db.Connect(db_username, db_password, ""); err != nil {
			panic(err)
		}
		defer db.Close()

		// Clean database
		if err := db.CreateDatabase(db_name); err != nil {
			panic(err)
		}
	*/

	if err := db.Connect(db_username, db_password, db_name); err != nil {
		panic(err)
	}

	if err := db.Run("./table_create.sql"); err != nil {
		panic(err)
	}

	//http.HandleFunc("/inves", getRoot)
	//http.HandleFunc("/hello", getHello)

	http.ListenAndServe(":8080", nil)
}
