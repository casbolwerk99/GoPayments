package main

import (
	"fmt"
	"net/http"

	"github.com/numeral/internal/payment"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	db, err := payment.InitializeDB()
	if err != nil {
		fmt.Println("Error initializing DB:", err)
		return
	}
	defer db.Close()

	handler := payment.NewHandler(db)

	http.HandleFunc("/payment-request", handler.HandleCreatePayment)
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error serving the server:", err)
		return
	}
}
