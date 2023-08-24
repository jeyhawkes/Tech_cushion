package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type KeyValueMap map[string]string

type Database struct {
	conn *sql.DB
}

// / with defaults
func (db *Database) ConnectDefault() error {
	return db.Connect("root", "password", "cushon")
}

func (db *Database) Connect(username, password, database string) error {
	var err error

	s := fmt.Sprintf("%s:%s@/%s?multiStatements=true", username, password, database)

	db.conn, err = sql.Open("mysql",
		s)

	return err
}

func (db *Database) SELECT(table string, keys string, where string, rows **sql.Rows) error {
	var query string

	if where == "" {
		query = fmt.Sprintf("SELECT %s FROM %s", keys, table)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s", keys, table, where)
	}

	var err error
	*rows, err = db.conn.Query(query)

	return err
}

func (db *Database) INSERT(table string, keyValuePairs KeyValueMap) error {

	var keys, values strings.Builder
	var first bool = true
	for k, v := range keyValuePairs {

		// stop comma from being at the end
		if first {
			first = false
		} else {
			keys.WriteByte(',')
			values.WriteByte(',')
		}

		keys.WriteString(k)
		values.WriteString(v)
	}

	query := fmt.Sprintf("INSERT %s (%s) VALUES (%s)", table, &keys, &values)
	var err error
	_, err = db.conn.Query(query)

	return err
}

func (db *Database) UPDATE(table string, keyValuePairs KeyValueMap, where string) error {

	var sets strings.Builder
	var first bool = true
	for k, v := range keyValuePairs {

		// stop comma from being at the end
		if first {
			first = false
		} else {
			sets.WriteByte(',')
		}

		sets.WriteString(k + "=" + v)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s;", table, &sets, where)
	var err error
	_, err = db.conn.Query(query)

	return err
}

func (db *Database) CountRows(table string, where string) (int, error) {

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s;", table, where)
	var rows *sql.Rows
	rows, err := db.conn.Query(query)

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		i += 1
	}

	return i, nil
}

func (db *Database) Run(path string) error {

	query, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if _, err := db.conn.Exec(string(query)); err != nil {
		return err
	}

	return nil
}

func (db *Database) CreateDatabase(name string) error {

	query := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", name)
	var err error
	_, err = db.conn.Query(query)

	if err != nil {
		return err
	}

	query = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`;", name)
	_, err = db.conn.Query(query)

	return err
}

func (db *Database) Close() error {
	return db.conn.Close()
}
