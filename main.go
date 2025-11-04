// main.go
package main

import (
	"database/sql"
	"db-migration-tool/internal/assert"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		slog.Error("Unable to open sql db for migrations", "err", err)
		os.Exit(1)
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
		slog.Error("Unable to create migration table", "err", err)
		os.Exit(1)
	}
	slog.Info("Connected! schema_migrations table ready.")

	files := getMigrationFiles()
	fileMap := buildFileMap(files)
	keys := make([]int, 0, len(fileMap))
	for k := range fileMap {
		keys = append(keys, k)
		sort.Ints(keys)
	}

	for _, version := range keys {
		if isApplied(db, version) {
			slog.Info("Version already applied", "version", version)
			continue
		}
		// fmt.Println(file)
		contents, err := getMigrationFileContents(fileMap[version])
		if err != nil {
			slog.Error("Unable to get contents for migration file", "err", err)
			os.Exit(1)
		}
		err = applyMigration(db, version, contents)
		if err != nil {
			slog.Error("Unable to apply migration", "err", err)
			os.Exit(1)
		}
		slog.Info("Migration applied", "version", version)
	}
}

func buildFileMap(files []string) map[int]string {
	fileMap := make(map[int]string)
	for _, file := range files {
		fileName := strings.SplitN(file, "/", 2)[1]
		versionNumber, err := strconv.Atoi(strings.SplitN(fileName, "_", 2)[0])
		if err != nil {
			slog.Error("file name in wrong format. expected 'migrations/<version>_<name>.sql'", "file_name", fileName)
			break
		}
		fileMap[versionNumber] = file
	}
	return fileMap
}

func getMigrationFiles() []string {
	files, err := filepath.Glob("./migrations/*up.sql")
	if err != nil {
		slog.Error("Unable to retrieve files for migration", "err", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		slog.Warn("No migration files found")
		os.Exit(0)
	}
	sort.Strings(files)
	return files
}

func getMigrationFileContents(filePath string) (string, error) {
	slog.Info("Filepath", "filepath", filePath)
	assert.Assert(
		strings.HasSuffix(filePath, ".sql"),
		"Migration file must be a .sql file",
	)
	assert.Assert(
		strings.HasPrefix(filePath, "migrations/"),
		"Migrations expected to be in ./migrations folder",
	)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	fmt.Printf("file contents: \n%s", string(data))
	return string(data), nil
}

func isApplied(db *sql.DB, v int) bool {
	slog.Info("Checking to see if version already exists in schema_migrations", "version", v)
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = (?)`, v).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func applyMigration(db *sql.DB, version int, migration string) error {
	slog.Info("Applying migration", "version", version)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(migration)
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
	if err != nil {
		return err
	}
	return tx.Commit()
}
