package main

import (
	"encoding/json"
	"log"
	"os"
)

func loadData() {
	if _, err := os.Stat("data/data.json"); err != nil {
		return // file doesn't exist yet, start fresh
	}

	file, err := os.Open("data/data.json")
	if err != nil {
		log.Printf("Error opening data file: %v", err)
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&app); err != nil {
		log.Printf("Error decoding data file: %v", err)
	}

	if app.Users == nil {
    app.Users = make(map[string]*User)
	}
	if app.Classrooms == nil {
		app.Classrooms = make(map[string]*Classroom)
	}
	if app.ShopItems == nil {
		app.ShopItems = make(map[string]*ShopItem)
	}
	if app.AvatarConfigs == nil {
		app.AvatarConfigs = make(map[string]*AvatarConfig)
	}
}

func saveData() {
	data, err := json.MarshalIndent(app, "", "    ")
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	err = os.WriteFile("data/data.json", data, 0644)
	if err != nil {
		log.Printf("Error writing data file: %v", err)
	}
}