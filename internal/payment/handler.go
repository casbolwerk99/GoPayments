package payment

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/santhosh-tekuri/jsonschema"
)

const (
	paymentSchemaPath = "data/request_schema.json"
	// in case of running tests, use the following path
	// paymentSchemaPath = "../../data/request_schema.json"
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
		panic(fmt.Sprintf("Error compiling JSONSchema: %v", err))
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

func (h *Handler) HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
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

	paymentInsertionChannel := make(chan error)

	go func() {
		err := InsertPayment(h.db, payment)
		if err != nil {
			paymentInsertionChannel <- fmt.Errorf("failed to insert payment: %w", err)
			return
		}

		paymentInsertionChannel <- nil
	}()

	go func() {
		if err := WritePaymentToBank(payment, os.Getenv("BANK_FOLDER")); err != nil {
			fmt.Println("Failed to write payment to bank", err)
		}
	}()

	go func() {
		if err := ProcessBankResponse(h.db, payment.IdempotencyUniqueKey); err != nil {
			fmt.Println("Failed to listen for and process bank response", err)
		}
	}()

	select {

	case err := <-paymentInsertionChannel:
		fmt.Println("Payment channel updated")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payment)

	case <-time.After(30 * time.Second):
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)

	}
}
