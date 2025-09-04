package database

import (
	"database/sql"
	_ "embed"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// Connection represents a function that returns a database connection.
// It provides lazy initialization of database connections,
type Connection func() (*sql.DB, error)

// schemes contains the embedded SQL schema definitions for database initialization.
//go:embed schemes.sql
var schemes string

// dbName defines the SQLite in-memory database connection string.
// :memory: creates a temporary database that exists only during execution
// _timeout=5000: Sets connection timeout to 5 seconds
const dbName = ":memory:?_timeout=5000&_sync=1"

// Shared connection variables
var (
	sharedDB  *sql.DB
	dbOnce    sync.Once
	dbInitErr error
)

// Initialize creates a Connection function that establishes a SQLite database
// connection in the system's temporary directory.
func Initialize() Connection {
	return func() (*sql.DB, error) {
		// Use sync.Once to ensure the database is only initialized once
		dbOnce.Do(func() {
			conn, err := sql.Open("sqlite3", dbName)
			if err != nil {
				dbInitErr = err
				return
			}
			sharedDB = conn
		})

		if dbInitErr != nil {
			return nil, dbInitErr
		}

		return sharedDB, nil
	}
}

// RunMigrations executes the embedded SQL schema against the database.
// It takes a Connection function to establish the database connection
// and applies all schema definitions from the embedded schemes.sql file.
func RunMigrations(conn Connection) error {
	db, err := conn()
	if err != nil {
		return err
	}

	if _, err := db.Exec(schemes); err != nil {
		return err
	}

	return nil
}
