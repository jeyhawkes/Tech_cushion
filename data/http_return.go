package data

import (
	"net/http"
)

type ErrorCode int32

type HTTPError struct {
	StatusCode int32
	ErrMessage string
}

// Use iota because the numbers don't map to anything (e.g. values in database)
type errorHTTPStatus int8

const (
	ErrorSuccess errorHTTPStatus = iota
	ErrorDatabaseRead
	ErrorDatabaseWrite
	ErrorJsonRead
	ErrorJsonWrite
	ErrorInvalidReqeust
	ErrorAlreadyExists
	ErrorUnknown
)

// Can't use http.StatusCode as a key because it would not allow errors to have the same error code
var ErrorHTTP = map[errorHTTPStatus]HTTPError{
	ErrorSuccess:        {http.StatusOK, ""},
	ErrorDatabaseRead:   {http.StatusInternalServerError, "error reading Database"},
	ErrorDatabaseWrite:  {http.StatusInternalServerError, "error writing Database"},
	ErrorJsonRead:       {http.StatusInternalServerError, "error reading json"},
	ErrorJsonWrite:      {http.StatusInternalServerError, "error writing json"},
	ErrorInvalidReqeust: {http.StatusMethodNotAllowed, "invalid method"},
	ErrorAlreadyExists:  {http.StatusConflict, "Already exists"},
	ErrorUnknown:        {http.StatusInternalServerError, "Unknown Error"},
}

type HTTPReturnData struct {
	Transaction_Id int       `json:"Transaction_Id"`
	Error_Code     ErrorCode `json:"Error_Code"`
	Timestamp      int64     `json:"Timestamp"`
	Error_Message  string    `json:"Error_Message"`
	Data           string    `json:"Data"`
}
