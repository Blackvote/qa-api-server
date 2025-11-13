package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"qa-api-service/handlers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "qa-api-service/migrations"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {

	resetDB := flag.Bool("reset-db", false, "Drop all tables by running goose down-to 0 before running migrations")
	flag.Parse()

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" {
		user = "qa_user"
	}
	if pass == "" {
		pass = "qa_password"
	}
	if dbName == "" {
		dbName = "qa_db"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		user, pass, dbName,
	)

	log.Printf("Starting QA API service...")

	log.Printf("Running migrations...")
	if err := runMigrations(dsn, *resetDB); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}
	log.Printf("Migrations completed")

	log.Printf("Connecting to database '%s' as user '%s'...", dbName, user)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	log.Printf("Connected to database")

	http.HandleFunc("/questions/", logLittle(handlers.MakeQuestionsHandler(db)))
	http.HandleFunc("/answers/", logLittle(handlers.MakeAnswersHandler(db)))

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(dsn string, reset bool) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	const migrationsDir = "./migrations"

	if reset {
		log.Println("Reset-db flag detected: running goose down-to 0...")
		if err := goose.DownTo(db, migrationsDir, 0); err != nil {
			return fmt.Errorf("failed to reset db: %w", err)
		}
		log.Println("Database reset complete, applying fresh migrations...")
	}

	return goose.Up(db, migrationsDir)
}

func logLittle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}
