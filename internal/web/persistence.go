package web

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dragon123098/Attendance-HackDay.git/internal/domain"
)

var app domain.AppState

func LoadData() {
	if _, err := os.Stat("data/data.json"); err != nil {
		ensureAppState()
		return
	}

	file, err := os.Open("data/data.json")
	if err != nil {
		log.Printf("Error opening data file: %v", err)
		ensureAppState()
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&app); err != nil {
		log.Printf("Error decoding data file: %v", err)
	}

	ensureAppState()
}

func ensureAppState() {
	if app.Users == nil {
		app.Users = make(map[string]*domain.User)
	}
	if app.Classrooms == nil {
		app.Classrooms = make(map[string]*domain.Classroom)
	}
	if app.ShopItems == nil {
		app.ShopItems = make(map[string]*domain.ShopItem)
	}
	if app.OwnedShopItems == nil {
		app.OwnedShopItems = make(map[string][]string)
	}
	if app.AvatarConfigs == nil {
		app.AvatarConfigs = make(map[string]*domain.AvatarConfig)
	}
	if app.ManualCoinAdjustments == nil {
		app.ManualCoinAdjustments = make(map[string]int)
	}
	if app.Transactions == nil {
		app.Transactions = []domain.CoinTransaction{}
	}
	if app.Attendance == nil {
		app.Attendance = []domain.AttendanceRecord{}
	}
	if app.Schedule == nil {
		app.Schedule = []domain.Schedule{}
	}
}

func saveData() {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("Error creating data directory: %v", err)
		return
	}

	data, err := json.MarshalIndent(app, "", "    ")
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	if err := os.WriteFile("data/data.json", data, 0644); err != nil {
		log.Printf("Error writing data file: %v", err)
	}
}
