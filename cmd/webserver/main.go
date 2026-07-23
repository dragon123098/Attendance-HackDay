package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

<<<<<<< HEAD
	"github.com/PeterGrunig/Attendance-HackDay/internal/integrations"
=======
	"github.com/joho/godotenv"
>>>>>>> main
	"github.com/PeterGrunig/Attendance-HackDay/internal/store"
	"github.com/PeterGrunig/Attendance-HackDay/internal/web"

	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	db, err := sql.Open("postgres", databaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

<<<<<<< HEAD
	storeOptions := []store.SQLStoreOption{}
	credentialCipher, err := integrations.NewAESGCMCredentialCipher(os.Getenv("INTEGRATION_CREDENTIAL_KEY"))
	if err != nil {
		log.Printf("integration credential storage disabled: %v", err)
	} else {
		storeOptions = append(storeOptions, store.WithCredentialCipher(credentialCipher))
	}

	log.Print("starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", web.NewRouter(store.NewSQLStore(db, storeOptions...))))
=======
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Printf("starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, web.NewRouter(store.NewSQLStore(db))))
>>>>>>> main
}

func databaseURL() string {
	if value := os.Getenv("DATABASE_URL"); value != "" {
		return value
	}

	return "postgres://attendance:Password123!@localhost:5433/attendancehackday?sslmode=disable"
}
