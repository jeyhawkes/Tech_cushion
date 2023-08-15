package setup

import (
	"path/filepath"

	"github.com/jeyhawkes/tech_cushion/database"
)

func Db(db *database.Database) error {
	if err := db.Connect(db_username, db_password, ""); err != nil {
		return err
	}

	// Clean database
	if err := db.CreateDatabase(db_name); err != nil {
		return err
	}

	if err := db.Connect(db_username, db_password, db_name); err != nil {
		return err
	}

	path := filepath.Join("github.com/jeyhawkes/Tech_cushion/", "setup", "table_create.sql")
	if err := db.Run(path); err != nil {
		return err
	}

	return nil
}
