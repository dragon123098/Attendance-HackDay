package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dragon123098/Attendance-HackDay.git/internal/store"
	"github.com/dragon123098/Attendance-HackDay.git/internal/web"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	db, err := sql.Open("sqlserver", databaseURL())
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

	return "server=localhost;user id=sa;password=Password123!;database=AttendanceHackday;encrypt=disable"
}
