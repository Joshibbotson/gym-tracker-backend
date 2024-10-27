package main

import (
	"fmt"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
)

func main() {
	db.ConnectDB()
	defer db.DisconnectDB()

	fmt.Println("Running application...")
	fmt.Println("Starting My Go Project")

}
