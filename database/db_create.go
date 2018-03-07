package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const db_port = 30257
const db_host = "localhost"
const db_owner = "root"
const db_user = "maxroach"
const db_consumer_geo = "consumer_geo"

func main() {
	// Connect to the database.
	url := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable", db_owner, db_host, db_port, db_consumer_geo)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
		fmt.Println("error connecting to the database: ", err)
	}
	defer db.Close()
	//fmt.Println("sql.Open finished!!")

	ddls := []string{
		"CREATE USER " + db_user,
		"CREATE DATABASE " + db_consumer_geo,
		"GRANT ALL ON DATABASE " + db_consumer_geo + " TO " + db_user,
	}
	for _, ddl := range ddls {
		if _, err := db.Exec(ddl); err != nil {
			fmt.Printf("failed!! :%v\n", ddl)
			log.Print(err)
		} else {
			fmt.Printf("finished: %v\n", ddl)
		}
	}

	objs := []struct {
		class      string
		name       string
		definition string
	}{
		{class: "TABLE", name: "location",
			definition: `(
				id INT,
				time TIMESTAMPTZ,
				lat NUMERIC,
				lng NUMERIC,
				CONSTRAINT "primary" PRIMARY KEY (id,time)
			)`},
	}
	for _, obj := range objs {
		ddl := fmt.Sprintf("CREATE %s IF NOT EXISTS %s %s", obj.class, obj.name, obj.definition)
		if _, err := db.Exec(ddl); err != nil {
			fmt.Printf("create %v %v failed!!\n", obj.class, obj.name)
			log.Print(err)
		} else {
			fmt.Printf("create %v %v finished :%v \n", obj.class, obj.name, obj.definition)
		}
	}

}
