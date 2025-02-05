package payment

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
)

const (
	paymentSchemaPath = "data/request_schema.json"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func isAuthorizedRequest(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != "Basic" {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		return false
	}

	credentials := string(decoded)
	return credentials == os.Getenv("username")+":"+os.Getenv("password")
}

func isValidRequest(payment Payment) bool {
	schema, err := jsonschema.NewCompiler().Compile(paymentSchemaPath)
	if err != nil {
		log.Fatal(err)
	}
	instance, err := json.Marshal(payment)
	if err != nil {
		return false
	}

	if err = schema.Validate(bytes.NewReader(instance)); err != nil {
		fmt.Println("Error validating payment:", err)
		return false
	}

	return true
}

func HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	if !isAuthorizedRequest(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !isValidRequest(payment) {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// err := InsertPayment(h.db, payment)
	// if err != nil {
	// 	http.Error(w, "Failed to insert payment", http.StatusInternalServerError)
	// 	return
	// }

	if err := WritePaymentToBank(payment, os.Getenv("BANK_FOLDER")); err != nil {
		http.Error(w, "Failed to write payment to bank", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}
