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
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}

	return db, nil
}

func PrintDB(db *sql.DB) {
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

	stmt, err := db.Prepare("insert into payments(id, status) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	id := payment.IdempotencyUniqueKey
	status := "REQUESTED"

	_, err = stmt.Exec(id, status)
	if err != nil {
		return err
	}

	PrintDB(db)

	return nil
}

func InsertData(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
