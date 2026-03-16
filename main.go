package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sushil-lumio/user-api/db"
	"github.com/sushil-lumio/user-api/handlers"
)

func main() {
	database, err := db.GetDB("users.db")
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer database.Close()

	handlers.DB = database

	// Public routes
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// Protected routes
	http.HandleFunc("/users", handlers.AuthMiddleware(handlers.UserHandler))

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
