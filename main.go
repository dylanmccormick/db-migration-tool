// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create schema_migrations table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version INTEGER PRIMARY KEY,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected! schema_migrations table ready.")
}
