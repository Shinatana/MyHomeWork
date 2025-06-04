package migrate

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(dbURL, migrationsPath string) error {
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Printf("Failed to create migrate instance: %v", err)
		return err
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Nothing to migrate")
			return nil
		}
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Printf("Migrations successfully migrated")
	return nil
}
