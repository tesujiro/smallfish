package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to the "bank" database.
	db, err := sql.Open("postgres", "postgresql://maxroach@localhost:30257/bank?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
		fmt.Println("error connecting to the database: ", err)
	}
	fmt.Println("sql.Open finished!!")

	/*
			// Create the "accounts" table.
			if _, err := db.Exec(
				"CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT)"); err != nil {
				fmt.Println("create table failed!!")
				log.Fatal(err)
			}
			fmt.Println("create table finished!!")

		// Insert two rows into the "accounts" table.
		if _, err := db.Exec(
			"INSERT INTO accounts (id, balance) VALUES (1, 1000.123), (2, 250)"); err != nil {
			log.Fatal(err)
		}
		fmt.Println("insert table finished!!")
	*/

	// Print out the balances.
	rows, err := db.Query("SELECT id, balance FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("Initial balances:")
	for rows.Next() {
		var id int
		var balance float64
		if err := rows.Scan(&id, &balance); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v %v\n", id, balance)
	}
}
