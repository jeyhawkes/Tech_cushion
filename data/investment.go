package data

import (
	"github.com/jeyhawkes/tech_cushon/database"
)

// order to save space
type InvestmentData struct {
	Id           database.UTINYINT
	Created_Date database.TIMESTAMP
	Name         database.TINYTEXT
}

type InvestmentListHTTP struct {
	Investment_Type_Id database.UTINYINT
	Name               database.TINYTEXT
}
