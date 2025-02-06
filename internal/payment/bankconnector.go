package payment

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

const (
	bankResponsePath = "data/bank_response.csv"
)

func ProcessBankResponse(db *sql.DB, id string) error {
	time.Sleep(5 * time.Second)
	fmt.Println("Bank response received. Processing...")

	file, err := os.Open(bankResponsePath)
	if err != nil {
		fmt.Println("Failed to open CSV file:", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		fmt.Println("Failed to read CSV header:", err)
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Failed to read CSV file:", err)
		return err
	}

	for _, record := range records {
		if len(record) != 2 {
			fmt.Println("Skipping invalid record:", record)
			continue
		}
		_, status := record[0], record[1]

		err := UpdatePayment(db, id, status)
		if err != nil {
			fmt.Println("Failed to update DB for ID:", id, err)
			continue
		}
	}

	return nil
}
