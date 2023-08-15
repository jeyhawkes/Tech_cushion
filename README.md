# Tech_cushion

Standalone application to prototype the senario given.

The database is DROPPED and re-created each time on run time so that the senario can be tested everytime

Database Tables
*the tables use an auto incrementing id as a primary key and idenfier (save space and time) - in practice this should probably move towards a UUID

Customer

go test -v ./...

the tests show how each api will work in practice and some of the scenarios i considered

RESTful 
* Should add some level of auth (basicAuth, OAuth, etc...)
* v1 (client can only choose 1 fund at a time)
  
Endpoints:

GET /invest/list/v1/ - Get list of funds 

GET /invest/customer/v1/*customer_id* - Get fund info about customer

POST /invest/customer/v1/*customer_id* - Add customer money to fund
{ 
    investment_type_id: int, 
    Amount:             int,
}

PATCH /invest/customer/v1/*customer_id* - Update customers amount 
{ 
    investment_type_id: int, 
    Amount:             int,
}


