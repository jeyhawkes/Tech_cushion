package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jeyhawkes/tech_cushion/data"
	"github.com/jeyhawkes/tech_cushion/database"
	"github.com/jeyhawkes/tech_cushion/logger"
	"github.com/jeyhawkes/tech_cushion/setup"
)

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got string, want data.HTTPReturnData) {

	var out data.HTTPReturnData
	err := json.Unmarshal([]byte(got), &out)
	if err != nil {
		t.Errorf("response body is wrong, couldn't parse to json : %s", got)
	}

	// make sure time was within the last min
	now := time.Now().Unix()
	if out.Timestamp < now-int64(time.Minute) {
		t.Errorf("response body is wrong, invalid timestamp (got %d, want %d)", out.Timestamp, time.Now().Unix())
	}

	if out.Timestamp > now+int64(time.Minute) {
		t.Errorf("response body is wrong, invalid timestamp (got %d, want %d)", out.Timestamp, time.Now().Unix())
	}

	t.Helper()
	out.Transaction_Id = 0 // unable to validiate
	out.Timestamp = 0      // already validated

	if !reflect.DeepEqual(out, want) {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func TestGetList(t *testing.T) {

	var log logger.Logger
	var err error
	if log, err = logger.NewLogger("log.txt"); err != nil {
		t.Errorf("Couldn't create log : %s", err)
	}

	var db database.Database
	if err := db.ConnectDefault(); err != nil {
		t.Errorf("database error %s ", err.Error())
	}
	defer db.Close()
	investhttp := NewInvestmentHTTP(&db, &log)

	request, _ := http.NewRequest(http.MethodGet, "/invest/list/v1/", nil)
	response := httptest.NewRecorder()

	investhttp.HandleInvestment(response, request)

	errHTTP := data.ErrorHTTP[data.ErrorSuccess]
	want := data.HTTPReturnData{
		Transaction_Id: 0,
		Error_Code:     data.ErrorCode(errHTTP.StatusCode),
		Error_Message:  errHTTP.ErrMessage,
		Data:           "[{\"Investment_Type_Id\":1,\"Name\":\"Cushon Equities Fund\"},{\"Investment_Type_Id\":2,\"Name\":\"Cushon Fixed income Fund\"}]",
		Timestamp:      0,
	}

	assertStatus(t, response.Code, int(errHTTP.StatusCode))
	assertResponseBody(t, response.Body.String(), want)
}

func TestCreateCustomerInvestment(t *testing.T) {

	// create anonymous test struct
	tests := []struct {
		name        string
		customer_id database.UMEDUIMINT
		data        data.CustomerInvestmentHTTP
		want        data.HTTPReturnData
	}{
		{
			name:        "Valid",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 1,
				Amount:             25000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorSuccess].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorSuccess].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
		{
			name:        "AlreadyExists",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 1,
				Amount:             25000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorAlreadyExists].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorAlreadyExists].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
		{
			name:        "InvalidInvestmentTypes",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 99,
				Amount:             25000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorInvalidReqeust].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorInvalidReqeust].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
	}

	var log logger.Logger
	var err error
	if log, err = logger.NewLogger("log.txt"); err != nil {
		t.Errorf("Couldn't create log : %s", err)
	}

	var db database.Database
	if err := setup.Db(&db); err != nil {
		t.Errorf("Couldn't create db : %s", err)
	}
	defer db.Close()
	investhttp := NewInvestmentHTTP(&db, &log)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			jParam, err := json.Marshal(test.data)
			if err != nil {
				t.Errorf("response body is wrong, couldn't parse to json")
			}

			request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/invest/customer/v1/%d", test.customer_id), strings.NewReader(string(jParam)))
			response := httptest.NewRecorder()

			investhttp.HandleCustomerInvestment(response, request)

			assertStatus(t, response.Code, int(test.want.Error_Code))
			assertResponseBody(t, response.Body.String(), test.want)
		})

	}

}

func TestGetCustomerInvestment(t *testing.T) {

	// create anonymous test struct
	tests := []struct {
		name        string
		customer_id database.UMEDUIMINT
		want        data.HTTPReturnData
	}{
		{
			name:        "Valid",
			customer_id: 1,
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorSuccess].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorSuccess].ErrMessage,
				Data:           "[{\"Customer_Id\":1,\"Investment_Type_Id\":1,\"Amount\":25000}]",
				Timestamp:      0,
			},
		},
		{
			name:        "InvalidCustomerId",
			customer_id: 99,
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorInvalidReqeust].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorInvalidReqeust].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
	}

	var log logger.Logger
	var err error
	if log, err = logger.NewLogger("log.txt"); err != nil {
		t.Errorf("Couldn't create log : %s", err)
	}

	var db database.Database
	if err := db.ConnectDefault(); err != nil {
		t.Errorf("database error %s ", err.Error())
	}
	defer db.Close()
	investhttp := NewInvestmentHTTP(&db, &log)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/invest/customer/v1/%d", test.customer_id), nil)
			response := httptest.NewRecorder()

			investhttp.HandleCustomerInvestment(response, request)

			assertStatus(t, response.Code, int(test.want.Error_Code))
			assertResponseBody(t, response.Body.String(), test.want)
		})

	}
}

func TestUpdateCustomerInvestment(t *testing.T) {
	// create anonymous test struct
	tests := []struct {
		name        string
		customer_id database.UMEDUIMINT
		data        data.CustomerInvestmentHTTP
		want        data.HTTPReturnData
	}{
		{
			name:        "SameFundSameAmount",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 1,
				Amount:             25000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorSuccess].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorSuccess].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
		{
			name:        "SameFundNewAmount",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 1,
				Amount:             20000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorSuccess].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorSuccess].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
		{
			name:        "NewFundSameAmount",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 2,
				Amount:             20000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorSuccess].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorSuccess].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
		{
			name:        "NewFundNewAmount",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 1,
				Amount:             20000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorSuccess].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorSuccess].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
		{
			name:        "InvalidFund",
			customer_id: 1,
			data: data.CustomerInvestmentHTTP{
				Investment_Type_Id: 99,
				Amount:             25000,
			},
			want: data.HTTPReturnData{
				Transaction_Id: 0,
				Error_Code:     data.ErrorCode(data.ErrorHTTP[data.ErrorInvalidReqeust].StatusCode),
				Error_Message:  data.ErrorHTTP[data.ErrorInvalidReqeust].ErrMessage,
				Data:           "",
				Timestamp:      0,
			},
		},
	}

	var log logger.Logger
	var err error
	if log, err = logger.NewLogger("log.txt"); err != nil {
		t.Errorf("Couldn't create log : %s", err)
	}

	var db database.Database
	if err := db.ConnectDefault(); err != nil {
		t.Errorf("database error %s ", err.Error())
	}
	defer db.Close()
	investhttp := NewInvestmentHTTP(&db, &log)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			jParam, err := json.Marshal(test.data)
			if err != nil {
				t.Errorf("response body is wrong, couldn't parse to json")
			}

			request, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/invest/customer/v1/%d", test.customer_id), strings.NewReader(string(jParam)))
			response := httptest.NewRecorder()

			investhttp.HandleCustomerInvestment(response, request)

			assertStatus(t, response.Code, int(test.want.Error_Code))
			assertResponseBody(t, response.Body.String(), test.want)
		})

	}
}
