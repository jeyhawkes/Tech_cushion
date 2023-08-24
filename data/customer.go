package data

import "github.com/jeyhawkes/tech_cushon/database"

// order to save space
type CustomerInvestmentData struct {
	// https://www.natwestgroup.com/ 19,000 customers
	// In a real world senario this would be an UUID
	Customer_Id        database.UMEDUIMINT
	Investment_Type_Id database.UTINYINT
	Amount             database.UMEDUIMINT
}

type CustomerInvestmentHTTP struct {
	Investment_Type_Id database.UTINYINT
	Amount             database.UMEDUIMINT
	Timestamp          int64
}
