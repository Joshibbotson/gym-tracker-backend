package main

import (
	"fmt"
	"net/http"

	"github.com/joshibbotson/gym-tracker-backend/internal/db"
)

func main() {
	db.ConnectDB()
	defer db.DisconnectDB()
	// http.HandleFunc("/todo", handler)
	// put in env variable.
	http.ListenAndServe(":8888", nil)

	fmt.Println("Running application...")
	fmt.Println("Starting My Go Project")

}
