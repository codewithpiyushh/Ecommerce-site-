package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"fashion/database"
	"fashion/handlers"

	_ "github.com/go-sql-driver/mysql"
)

// Local DB connection (used only in main.go)
var con *sql.DB

// CORS middleware
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Step 1: Initialize DB connection
	err := database.InitConnection()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	con = database.Con

	// Step 2: Set handlers package's DB connection
	handlers.Con = database.Con
	// Step 3: Define routes
	http.Handle("/signup.html", enableCors(http.HandlerFunc(handlers.Servesignup)))
	http.Handle("/signup-submit", enableCors(http.HandlerFunc(handlers.SignupHandler)))
	http.Handle("/login.html", enableCors(http.HandlerFunc(handlers.ServeLogin)))
	http.Handle("/login-submit", enableCors(http.HandlerFunc(handlers.LoginHandler)))
	http.Handle("/", enableCors(http.HandlerFunc(handlers.IndexHandler)))
	http.Handle("/home_furnishing.html", enableCors(http.HandlerFunc(handlers.HomelivingHandler)))
	http.Handle("/mens.html", enableCors(http.HandlerFunc(handlers.MenHandler)))
	http.Handle("/women.html", enableCors(http.HandlerFunc(handlers.WomenHandler)))
	http.Handle("/cart.html", enableCors(handlers.JwtMiddleware(handlers.CartHandler)))
	http.Handle("/auth-check", enableCors(handlers.JwtMiddleware(http.HandlerFunc(handlers.AuthCheckHandler))))
	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/dashboard.html", handlers.JwtMiddleware(handlers.DashboardHandler))
	http.HandleFunc("/logout", handlers.LogoutHandler)

	http.HandleFunc("/cart", handlers.AddToCartHandler)
	http.HandleFunc("/displaycart", handlers.DisplayCartItemsHandler)
	http.HandleFunc("/deletecart", handlers.DeleteCartItemHandler)
	



	// Step 4: Start server
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
