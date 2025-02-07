package main

import (
	"fmt"
	"net/http"

	"github.com/numeral/internal/payment"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	db, err := payment.InitializeDB()
	if err != nil {
		fmt.Println("Error initializing DB")
		panic(fmt.Sprintf("Error initializing DB: %v", err))
	}
	defer db.Close()

	handler := payment.NewHandler(db)

	http.HandleFunc("/payment-request", handler.HandleCreatePayment)
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	http.ListenAndServe(port, nil)
}
