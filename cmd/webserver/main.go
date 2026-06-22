package main

import (
	"log"
	"net/http"

	"github.com/dragon123098/Attendance-HackDay.git/internal/web"
)

func main() {
	web.LoadData()

	log.Print("starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", web.NewRouter()))
}
