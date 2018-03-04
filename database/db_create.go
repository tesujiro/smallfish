package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const db_port = 30257
const db_host = "localhost"
const db_user = "root"
const db_consumer_geo = "consumer_geo"

func main() {
	// Connect to the database.
	url := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable", db_user, db_host, db_port, db_consumer_geo)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
		fmt.Println("error connecting to the database: ", err)
	}
	defer db.Close()
	fmt.Println("sql.Open finished!!")

	objs := []struct {
		class      string
		name       string
		definition string
	}{
		{class: "DATABASE", name: db_consumer_geo, definition: ""},
		{class: "TABLE", name: "location", definition: "(id INT PRIMARY KEY, lat NUMERIC, lng NUMERIC)"},
	}
	for _, obj := range objs {
		ddl := fmt.Sprintf("CREATE %s IF NOT EXISTS %s %s", obj.class, obj.name, obj.definition)
		if _, err := db.Exec(ddl); err != nil {
			fmt.Printf("create %v %v failed!!\n", obj.class, obj.name)
			log.Fatal(err)
		}
		fmt.Printf("create %v %v finished :%v \n", obj.class, obj.name, obj.definition)
	}
}
