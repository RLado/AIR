package setup

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Embed DB skeleton
//
//go:embed skel.sql
var skel string

// Creates the necessary database for the service in the given path
func CreateDatabase(db_path string) error {
	// If the database file already exists delete it
	if _, err := os.Stat(db_path); err == nil {
		err := os.Remove(db_path)
		if err != nil {
			return fmt.Errorf("error creating a new database, failed existing database file deletion: %s", err)
		}
	}

	// Open database and create a database file
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", db_path))
	if err != nil {
		return err
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(skel)
	if err != nil {
		return fmt.Errorf("error creating database tables: %s", err)
	}

	return nil
}
