package payment

type Payment struct {
	DebtorIban           string  `json:"debtor_iban"`
	DebtorName           string  `json:"debtor_name"`
	CreditorIban         string  `json:"creditor_iban"`
	CreditorName         string  `json:"creditor_name"`
	Ammount              float64 `json:"ammount"`
	IdempotencyUniqueKey string  `json:"idempotency_unique_key"`
}
