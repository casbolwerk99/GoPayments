package payment

import (
	"encoding/json"
	"fmt"
)

type Payment struct {
	DebtorIban           string `json:"debtor_iban"`
	DebtorName           string `json:"debtor_name"`
	CreditorIban         string `json:"creditor_iban"`
	CreditorName         string `json:"creditor_name"`
	Ammount              int64  `json:"ammount"` // in cents
	IdempotencyUniqueKey string `json:"idempotency_unique_key"`
}

func (payment *Payment) UnmarshalJSON(data []byte) error {
	// create an alias to avoid infinite recursion
	type Alias Payment

	var raw struct {
		Ammount float64 `json:"ammount"`
		*Alias          // embed the original struct to auto-populate all other fields
	}

	raw.Alias = (*Alias)(payment)

	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println("Error unmarshalling payment:", err)
		return err
	}

	payment.Ammount = int64(raw.Ammount * 100)
	return nil
}
