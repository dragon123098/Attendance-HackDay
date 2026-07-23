package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/PeterGrunig/Attendance-HackDay/internal/integrations"
	"github.com/PeterGrunig/Attendance-HackDay/internal/store"
	"github.com/PeterGrunig/Attendance-HackDay/internal/web"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", databaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storeOptions := []store.SQLStoreOption{}
	credentialCipher, err := integrations.NewAESGCMCredentialCipher(os.Getenv("INTEGRATION_CREDENTIAL_KEY"))
	if err != nil {
		log.Printf("integration credential storage disabled: %v", err)
	} else {
		storeOptions = append(storeOptions, store.WithCredentialCipher(credentialCipher))
	}

	log.Print("starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", web.NewRouter(store.NewSQLStore(db, storeOptions...))))
}

func databaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}

	return "postgres://attendance:Password123!@localhost:5433/attendancehackday?sslmode=disable"
}
