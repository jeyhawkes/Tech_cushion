package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jeyhawkes/tech_cushion/data"
	"github.com/jeyhawkes/tech_cushion/database"
	"github.com/jeyhawkes/tech_cushion/logger"
)

const ERROR_DATABASE_READ = "Could not read from data base"

// MUST:: Make sure databse and logger are passed by pointer so it can be run as go routine
type investmentHTTP struct {
	db            *database.Database
	log           *logger.Logger
	transactionId int
}

func NewInvestmentHTTP(db *database.Database, log *logger.Logger) investmentHTTP {
	return investmentHTTP{db: db, log: log}
}

func (invest_http *investmentHTTP) HandleInvestment(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	invest_http.transactionId += 1

	if req.Method == http.MethodGet {
		invest_http.GetInvestmentList(w, req, invest_http.transactionId)
		return
	}

	invest_http.writeHTTPOutput(w, invest_http.transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
}

func (invest_http *investmentHTTP) HandleCustomerInvestment(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// increament transactionId to keep track and log
	invest_http.transactionId += 1

	url_split := strings.Split(req.URL.Path, "/")

	invest_http.log.LogInfo(invest_http.transactionId, req.URL.Path)

	// Endpoints only require customer_id
	if len(url_split) != 5 {
		invest_http.log.LogError(invest_http.transactionId, "Invalid URL")
		invest_http.writeHTTPOutput(w, invest_http.transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
		return
	}

	customer_id, err := strconv.Atoi(url_split[4])

	if err != nil {
		invest_http.log.LogError(invest_http.transactionId, "Invalid Customer Id")
		invest_http.writeHTTPOutput(w, invest_http.transactionId, data.ErrorHTTP[data.ErrorAlreadyExists], "")
		return
	}

	switch req.Method {
	case http.MethodGet:
		invest_http.GetCustomerInvestment(w, req, invest_http.transactionId, database.UMEDUIMINT(customer_id))
		return
	case http.MethodPost:
		invest_http.CreateCustomerInvestment(w, req, invest_http.transactionId, database.UMEDUIMINT(customer_id))
		return
	case http.MethodPatch:
		invest_http.UpdateCustomerInvestment(w, req, invest_http.transactionId, database.UMEDUIMINT(customer_id))
		return
	}

	invest_http.writeHTTPOutput(w, invest_http.transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
}

func (invest_http *investmentHTTP) GetInvestmentList(w http.ResponseWriter, req *http.Request, transactionId int) {

	var list []data.InvestmentList
	if err := invest_http.getInvestmentList(&list); err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseRead], "")
		return
	}

	j, err := json.Marshal(list)
	if err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorJsonWrite], "")
		return
	}

	invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorSuccess], string(j))
}

/*
V1 CURRENT (client can only choose 1 fund at a time) :
- Already returns a list of investments for the customer (only 1 will return)

V2 (clients can have muliple funds) : NO CHANGES NEED
- Will automatically scale to allow the customer to have multiple investments
*/

func (invest_http *investmentHTTP) GetCustomerInvestment(w http.ResponseWriter, req *http.Request, transactionId int, customerId database.UMEDUIMINT) {

	invest_http.log.LogInfo(invest_http.transactionId, fmt.Sprintf("customer_id : %d", customerId))

	// is valid inputs
	exists, err := invest_http.validateInputs(&customerId, nil)
	if err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseRead], "")
		return
	} else if !exists {
		invest_http.log.LogError(transactionId, "Invalid input")
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
		return
	}

	var list []data.CustomerInvestmentData
	if err := invest_http.getCustomerInvestment(customerId, &list); err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseRead], "")
		return
	}

	j, err := json.Marshal(list)
	if err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorJsonWrite], "")
		return
	}

	invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorSuccess], string(j))
}

/*
V1 CURRENT (client can only choose 1 fund at a time) :
- If customer investment already exists, fail

V2 (clients can have muliple funds) :
- If customer investment already exists but NOT assoicated with the fund, pass
- If customer investment already exists but IS assoicated with the fund, fail
*/
func (invest_http *investmentHTTP) CreateCustomerInvestment(w http.ResponseWriter, req *http.Request, transactionId int, customerId database.UMEDUIMINT) {
	invest_http.log.LogInfo(invest_http.transactionId, fmt.Sprintf("customer_id : %d", customerId))
	invest_http.log.LogInfo(invest_http.transactionId, fmt.Sprintf("inputs : %d", req.Body))

	// Get post parameters
	decoder := json.NewDecoder(req.Body)
	var post_params data.CustomerInvestmentHTTP
	err := decoder.Decode(&post_params)

	if err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
	}

	// is valid inputs
	exists, err := invest_http.validateInputs(&customerId, &post_params.Investment_Type_Id)
	if err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseRead], "")
		return
	} else if !exists {
		invest_http.log.LogError(transactionId, "Invalid input")
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
		return
	}

	// check if already exists
	var rows *sql.Rows
	if invest_http.db.SELECT("customer_investments", "*", fmt.Sprintf("customer_id = %d", customerId), &rows) == nil {
		defer rows.Close()
		if rows.Next() {
			invest_http.log.LogError(transactionId, "Already Exists")
			invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorAlreadyExists], "")
			return
		}
	}

	var db_params database.KeyValueMap = database.KeyValueMap{
		"customer_id":        fmt.Sprintf("%d", customerId),
		"investment_type_id": fmt.Sprintf("%d", post_params.Investment_Type_Id),
		"amount":             fmt.Sprintf("%d", post_params.Amount),
	}

	if err := invest_http.db.INSERT("customer_investments", db_params); err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseWrite], "")
		return
	}

	invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorSuccess], "")
}

/*
V1 CURRENT (client can only choose 1 fund at a time) :
- Uses customer_id as the WHERE as they can have only row. (less calls)
- Can update both fund and amount at once

V2 (clients can have muliple funds) :
- Will have to make an additional call to get the ID of the row
*/

func (invest_http *investmentHTTP) UpdateCustomerInvestment(w http.ResponseWriter, req *http.Request, transactionId int, customerId database.UMEDUIMINT) {
	invest_http.log.LogInfo(invest_http.transactionId, fmt.Sprintf("customer_id : %d", customerId))
	invest_http.log.LogInfo(invest_http.transactionId, fmt.Sprintf("inputs : %d", req.Body))

	// Get post parameters
	decoder := json.NewDecoder(req.Body)
	var post_params data.CustomerInvestmentHTTP
	err := decoder.Decode(&post_params)

	if err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
	}

	// is valid inputs
	exists, err := invest_http.validateInputs(&customerId, &post_params.Investment_Type_Id)
	if err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseRead], "")
		return
	} else if !exists {
		invest_http.log.LogError(transactionId, "Invalid input")
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
		return
	}

	// check if already exists
	var rows *sql.Rows
	if err = invest_http.db.SELECT("customer_investments", "*", fmt.Sprintf("customer_id = %d", customerId), &rows); err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseRead], "")
	}

	defer rows.Close()
	if !rows.Next() {
		invest_http.log.LogError(transactionId, "Doesn't exist")
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorInvalidReqeust], "")
		return
	}

	var db_params database.KeyValueMap = database.KeyValueMap{
		"investment_type_id": fmt.Sprintf("%d", post_params.Investment_Type_Id),
		"amount":             fmt.Sprintf("%d", post_params.Amount),
	}

	if err := invest_http.db.UPDATE("customer_investments", db_params, fmt.Sprintf("customer_id = %d", customerId)); err != nil {
		invest_http.log.LogError(transactionId, err.Error())
		invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorDatabaseWrite], "")
		return
	}

	invest_http.writeHTTPOutput(w, transactionId, data.ErrorHTTP[data.ErrorSuccess], "")
}

/// :: PRIVATE :: ///

func (invest_http *investmentHTTP) writeHTTPOutput(w http.ResponseWriter, transactionId int, errHTTP data.HTTPError, jout string) {

	out := data.HTTPReturnData{
		Transaction_Id: transactionId,
		Error_Code:     data.ErrorCode(errHTTP.StatusCode),
		Error_Message:  errHTTP.ErrMessage,
		Data:           jout,
		Timestamp:      time.Now().Unix(),
	}

	jData, err := json.Marshal(out)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(int(errHTTP.StatusCode))
	w.Write(jData)
}

func (invest_http *investmentHTTP) getInvestmentList(list *[]data.InvestmentList) error {

	var rows *sql.Rows
	if err := invest_http.db.SELECT("investment_types", "id, name", "", &rows); err != nil {
		return err
	}

	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var investment data.InvestmentList
		if err := rows.Scan(&investment.Investment_Type_Id, &investment.Name); err == nil {
			*list = append(*list, investment)
		} else {
			return err
		}
	}

	return nil
}

func (invest_http *investmentHTTP) validateInputs(customerId *database.UMEDUIMINT, investmentTypeId *database.UTINYINT) (bool, error) {

	var exists bool
	var err error

	if customerId != nil {
		exists, err = invest_http.isValidId("customer", int64(*customerId))
		if err != nil {
			return false, err
		} else if !exists {
			return false, nil
		}
	}

	if investmentTypeId != nil {
		exists, err = invest_http.isValidId("investment_types", int64(*investmentTypeId))
		if err != nil {
			return false, err
		} else if !exists {
			return false, nil
		}
	}

	return true, nil
}

func (invest_http *investmentHTTP) isValidId(table string, id int64) (bool, error) {

	// id exists
	i, err := invest_http.db.CountRows(table, fmt.Sprintf("id = %d", id))
	if err != nil {
		return false, err
	} else if i == 0 {
		return false, nil
	}

	return true, nil
}

// Return list so it can be scaled up to allow multiple investments
func (invest_http *investmentHTTP) getCustomerInvestment(customer_id database.UMEDUIMINT, list *[]data.CustomerInvestmentData) error {
	var rows *sql.Rows

	if err := invest_http.db.SELECT("customer_investments", "customer_id, investment_type_id, amount", fmt.Sprintf("customer_id = %d", customer_id), &rows); err != nil {
		return err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var customer_investment data.CustomerInvestmentData
		if err := rows.Scan(&customer_investment.Customer_Id, &customer_investment.Investment_Type_Id, &customer_investment.Amount); err == nil {
			*list = append(*list, customer_investment)
		} else {
			return err
		}
	}

	return nil
}
