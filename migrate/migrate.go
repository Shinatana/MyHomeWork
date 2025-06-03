package migrate

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"log"
)

func RunMigrations(dbURL, migrationsPath string) error {
	m, err := migrate.New(
		"file:///Users/konoko/Documents/Go/Homework3/study/migrations"+migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Printf("Failed to create migrate instance: %v", err)
		return err
		//уточнить
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Nothing to migrate")
			return nil
		}
		log.Printf("Failed to migrate: %v", err)
	}
	log.Printf("Migrations successfully migrated")
	return nil
}
