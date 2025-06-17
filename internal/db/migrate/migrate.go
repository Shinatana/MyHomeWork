package migrate

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"strings"
)

func RunMigrations(dbURL string) error {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter migration direction (up/down): ")
	directionInput, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read migration direction: %v", err)
	}
	directionInput = strings.TrimSpace(directionInput)

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
		"file://sql/migrations",
		"pgx",
		driver)
	if err != nil {
		return err
	}

	switch directionInput {
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
		return fmt.Errorf("unknown migration direction: %s", directionInput)
	}
	return nil
}
