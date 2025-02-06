// based on https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go

package payment

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB() (*sql.DB, error) {
	os.Remove(os.Getenv("SQLITE_DB_FILE_LOCATION"))

	fmt.Println("Initializing DB...")

	db, err := sql.Open("sqlite3", os.Getenv("SQLITE_DB_FILE_LOCATION"))
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	create table payments (id text not null primary key, status text);
	delete from payments;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println(err, ":", sqlStmt)
		return nil, err
	}

	return db, nil
}

func PrintDB(db *sql.DB) {
	fmt.Println("Printing DB...")
	rows, err := db.Query("select id, status from payments")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var status string
		err = rows.Scan(&id, &status)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, status)
	}
}

func InsertPayment(db *sql.DB, payment Payment) error {
	// inspired by https://github.com/mattn/go-sqlite3/blob/master/_example/json/json.go

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Failed to start transaction:", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into payments(id, status) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	id := payment.IdempotencyUniqueKey
	status := "REQUESTED"

	fmt.Println("Inserting payment with id:", id)

	_, err = stmt.Exec(id, status)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Failed to commit transaction:", err)
		return err
	}

	PrintDB(db)

	return nil
}

func UpdatePayment(db *sql.DB, id string, status string) error {
	// inspired by https://github.com/mattn/go-sqlite3/blob/master/_example/json/json.go

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Failed to start transaction:", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`update payments set status = ? where id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Println("Updating payment with id:", id)

	_, err = stmt.Exec(status, id)
	if err != nil {
		fmt.Println("Failed to execute SQL statement:", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Failed to commit transaction:", err)
		return err
	}

	PrintDB(db)

	return nil
}
