package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/jackc/pgx/v5"
	"log"
)

func RunMigrations(direction *string, dbURL string) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"pgx",
		driver)
	if err != nil {
		return err
	}

	switch *direction {
	case "up":
		err = m.Up()
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No new migrations to apply, database is up to date")
				return nil
			} else {
				log.Printf("migration failed: %v", err)
				return err
			}
		}
		log.Println("Migrations up completed")
	case "down":
		err = m.Down()
		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No new migrations to apply, database is up to date")
				return nil
			} else {
				log.Printf("migration failed: %v", err)
				return err
			}
		}
		log.Println("Migrations down completed")
	default:
		return fmt.Errorf("unknown migration direction: %s", direction)
	}
	return nil
}
