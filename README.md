# Tech_cushion

Standalone application to prototype the senario given.

The database is DROPPED and re-created each time on run time so that the senario can be tested everytime

Database Tables
*the tables use an auto incrementing id as a primary key and idenfier (save space and time) - in practice this should probably move towards a UUID

Customer

go test ./.. or
go test ./handlers/

the tests show how each api will work in practice and some of the scenarios i considered

RESTful 
Endpoints:
POST /invest/customer/v1/*customer_id*

GET /invest/customer/v1/*customer_id*

PATCH /invest/customer/v1/*customer_id*

GET /invest/list/v1/

