package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dragon123098/Attendance-HackDay.git/internal/store"
	"github.com/dragon123098/Attendance-HackDay.git/internal/web"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", databaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Print("starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", web.NewRouter(store.NewSQLStore(db))))
}

func databaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}

	return "postgres://attendance:Password123!@localhost:5433/attendancehackday?sslmode=disable"
}
