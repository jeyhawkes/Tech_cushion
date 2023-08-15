# Tech_cushion

## Standalone application to prototype the senario given.
The database is DROPPED and re-created each time on run time so that the senario can be tested everytime

## Database Tables
* The tables use an auto incrementing id as a primary key and idenfier (save space and time) - in practice this should probably move towards a UUID

```SQL
-- structure for table cushion.customer
CREATE TABLE `customer` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Allows almost 17 million and natwest has 19million customers so might have to be updated (still use UUID in pratice)',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'customer name',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created timestamp',
  PRIMARY KEY (`id`)
)
```

```SQL
-- structure for table cushion.investment_types
CREATE TABLE `investment_types` (
  `id` tinyint(3) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Small because there are probably not more than 255 different fund (Would still move to UUID in pratice)', 
  `name` varchar(255) NOT NULL DEFAULT '',
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)
```

```SQL
CREATE TABLE `customer_investments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Not used currently (v2 will become important to use to update indiviual customer funds so still move to UUID)  '
  `investment_type_id` tinyint(3) unsigned NOT NULL,
  `customer_id` mediumint(8) unsigned NOT NULL,
  `amount` mediumint(8) unsigned NOT NULL,
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)
```

## Tests 
The tests show how each api will work in practice and some of the scenarios i considered

go test -v ./...

## RESTful 
* v1 (client can only choose 1 fund at a time)
* v2 (clients can have muliple funds) - NOT IMPLEMENT (just nots for how this would be developed)
* Should add some level of auth (basicAuth, OAuth, etc...)
* Should add a timestamp on incoming comms (protection against replay attacks)
  
## Endpoints:

### GET /invest/list/v1/ - Get list of funds 

### GET /invest/customer/v1/*customer_id* - Get fund info about customer
v1 CURRENT: Already returns a list of investments for the customer (only 1 will return)

v2 NEXT STEP: NO CHANGES NEED: Will automatically scale to allow the customer to have multiple investments

### POST /invest/customer/v1/*customer_id* - Add customer money to fund
```golang
{ 
    investment_type_id: int
    amount:             int
}
```
v1 CURRENT: 
- If customer investment already exists, fail

v2 NEXT STEP:
- If customer investment already exists but NOT assoicated with the fund, pass
- If customer investment already exists but IS assoicated with the fund, fail


### PATCH /invest/customer/v1/*customer_id* - Update customers amount 
```golang
{ 
    investment_type_id: int 
    amount:             int
}
```
v1 CURRENT (client can only choose 1 fund at a time) :
- Uses customer_id as the WHERE as they can have only row. (less calls)
- Can update both fund and amount at once

v2 NEXT STEP: : 
- Will have to make an additional call to get the ID of the row

## Endpoints return 
```golang
{
	Transaction_Id int       `json:"Transaction_Id"`
	  Error_Code     ErrorCode `json:"Error_Code"`
	  Timestamp      int64     `json:"Timestamp"`
	  Error_Message  string    `json:"Error_Message"`
	  Data           string    `json:"Data"`
}
```
