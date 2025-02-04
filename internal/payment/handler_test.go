package payment

import (
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
)

func setTestEnvVars() {
	os.Setenv("username", "test")
	os.Setenv("password", "test")
}

func TestIsAuthorizedRequest(t *testing.T) {
	setTestEnvVars()

	tests := []struct {
		name           string
		authHeader     string
		expectedResult bool
	}{
		{
			name:           "Valid Basic Auth",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("username"), os.Getenv("password")))),
			expectedResult: true,
		},
		{
			name:           "Invalid Basic Auth",
			authHeader:     "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("username"), "wrongpassword"))),
			expectedResult: false,
		},
		{
			name:           "No Authorization Header",
			authHeader:     "",
			expectedResult: false,
		},
		{
			name:           "Invalid Authorization Header Format",
			authHeader:     "Bearer abcdefghijklmnopqrstuvwxyz",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/payment-request", nil)
			req.Header.Set("Authorization", tt.authHeader)

			authenticated := isAuthorizedRequest(req)
			if authenticated != tt.expectedResult {
				t.Errorf("Expected %v, but got %v", tt.expectedResult, authenticated)
			}
		})
	}
}

func TestIsValidRequest(t *testing.T) {
	setTestEnvVars()

	tests := []struct {
		name          string
		payment       Payment
		expectedValid bool
	}{
		{
			name: "Valid Payment",
			payment: Payment{
				DebtorIban:           "FR1112739000504482744411A64",
				DebtorName:           "company1",
				CreditorIban:         "DE65500105179799248552",
				CreditorName:         "beneficiary",
				Ammount:              42.99,
				IdempotencyUniqueKey: "JXJ984XXXZ",
			},
			expectedValid: true,
		},
		{
			name: "Invalid DebtorIban",
			payment: Payment{
				DebtorIban:           "123",
				DebtorName:           "company1",
				CreditorIban:         "DE65500105179799248552",
				CreditorName:         "beneficiary",
				Ammount:              42.99,
				IdempotencyUniqueKey: "JXJ984XXXZ",
			},
			expectedValid: false,
		},
		{
			name: "Invalid DebtorName",
			payment: Payment{
				DebtorIban:           "FR1112739000504482744411A64",
				DebtorName:           "a",
				CreditorIban:         "DE65500105179799248552",
				CreditorName:         "beneficiary",
				Ammount:              42.99,
				IdempotencyUniqueKey: "JXJ984XXXZ",
			},
			expectedValid: false,
		},
		{
			name: "Missing DebtorIban",
			payment: Payment{
				DebtorName:           "company1",
				CreditorIban:         "DE65500105179799248552",
				CreditorName:         "beneficiary",
				Ammount:              42.99,
				IdempotencyUniqueKey: "JXJ984XXXZ",
			},
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidRequest(tt.payment)
			if result != tt.expectedValid {
				t.Errorf("Expected %v, but got %v", tt.expectedValid, result)
			}
		})
	}
}
