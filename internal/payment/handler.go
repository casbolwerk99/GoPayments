package payment

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var (
	// TODO: make this relative and not OS specific
	paymentSchemaPath = "file:///D:/Documents/go/numeral/data/request_schema.json"
)

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
	schemaLoader := gojsonschema.NewReferenceLoader(paymentSchemaPath)
	documentLoader := gojsonschema.NewGoLoader(payment)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Println("Error validating payment:", err)
		return false
	}

	if !result.Valid() {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
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

	bankfolder := os.Getenv("BANK_FOLDER")
	if err := WritePaymentToBank(payment, bankfolder); err != nil {
		http.Error(w, "Failed to write payment to bank", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}
