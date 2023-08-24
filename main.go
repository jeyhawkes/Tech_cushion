package main

import (
	"net/http"

	"github.com/jeyhawkes/tech_cushon/database"
	"github.com/jeyhawkes/tech_cushon/handlers"
	"github.com/jeyhawkes/tech_cushon/logger"
	"github.com/jeyhawkes/tech_cushon/setup"
)

func main() {

	var log logger.Logger
	var err error
	if log, err = logger.NewLogger("log.txt"); err != nil {
		panic(err)
	}

	var db database.Database

	if err := setup.Db(&db); err != nil {
		panic(err)
	}
	defer db.Close()

	investhttp := handlers.NewInvestmentHTTP(&db, &log)
	http.HandleFunc("/invest/list/v1/", investhttp.HandleInvestment)
	http.HandleFunc("/invest/customer/v1/", investhttp.HandleCustomerInvestment)

	http.ListenAndServe(":8080", nil)
}
