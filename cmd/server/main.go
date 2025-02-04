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
		return
	}

	http.HandleFunc("/payment-request", payment.HandleCreatePayment)
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	http.ListenAndServe(port, nil)
}
