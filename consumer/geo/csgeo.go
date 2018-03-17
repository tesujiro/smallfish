package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Sleeper(w http.ResponseWriter, r *http.Request) {
	log.Printf("Sleeper!!")
	q := r.URL.Query()
	i, _ := strconv.Atoi(q.Get("timer"))
	time.Sleep(time.Duration(i) * time.Millisecond)
	fmt.Fprintf(w, "Hello, World :slept %d msec\n", i)
}

const db_user = "maxroach"
const db_consumer_geo = "consumer_geo"
const db_host_default = "localhost"
const db_port_default = 30257

//const db_host = "cockroachdb-public"
var db_host string
var db_port int

func connect() (*sql.DB, error) {
	if db_host == "" {
		db_host = db_host_default
	}
	if db_port == 0 {
		db_port = db_port_default
	}
	url := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable", db_user, db_host, db_port, db_consumer_geo)
	db, err := sql.Open("postgres", url)
	if err != nil {
		//log.Fatal("error connecting to the database: ", err)
		fmt.Println("error connecting to the database: ", err)
	}
	return db, err
}

type ConsumerGeoInfo struct {
	ConsumerId int     `json:"consumerId"`
	Lat        float64 `json:"latitude"`
	Lng        float64 `json:"longtitude"`
}

func addConsumerGeo(db *sql.DB, geo ConsumerGeoInfo) error {
	// Insert two rows into the "location" table.
	//stmt, err := db.Prepare("INSERT INTO location (id, time, lat, lng) VALUES (?,?,?,?)")
	stmt, err := db.Prepare("INSERT INTO location (id, time, lat, lng) VALUES ($1,now(),$2,$3)")
	if err != nil {
		log.Printf("prepare statement faled!!")
		return err
	}
	defer stmt.Close() // danger!
	log.Printf("prepare statement!!")

	//res, err := stmt.Exec(geo.ConsumerId, time.Now(), geo.Lat, geo.Lng)
	res, err := stmt.Exec(geo.ConsumerId, geo.Lat, geo.Lng)
	if err != nil {
		log.Printf("exec statement faled!!")
		return err
	}
	log.Printf("execute statement!!")

	/*
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	*/
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("get rows affected faled!!")
		return err
	}
	log.Printf("affected = %d\n", rowCnt)

	return nil
}

func ConsumerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConsumerHandler!!")
	vars := mux.Vars(r)
	lat, err := strconv.ParseFloat(vars["latitude"], 64)
	if err != nil {
		log.Fatal(err)
		return
	}
	lng, err := strconv.ParseFloat(vars["longtitude"], 64)
	if err != nil {
		log.Fatal(err)
		return
	}

	geo := ConsumerGeoInfo{Lat: lat, Lng: lng}

	db, err := connect()
	if err != nil {
		log.Printf("database connect failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("connected database!!")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("transaction begin failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	log.Printf("transaction begin!!")

	if err := addConsumerGeo(db, geo); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("commit faled!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("commit finished!!")

	fmt.Println("insert table finished!!")
	log.Printf("geo=%v\n", geo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusNotFound)

	if err := json.NewEncoder(w).Encode(geo); err != nil {
		log.Print("json.NewEncoder error!\n")
		log.Fatal(err)
	}
}

func ConsumerGeoCollectionWriter(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConsumerGeoCollectionWriter!!")

	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("bad Content-Type!!")
		log.Printf(r.Header.Get("Content-Type"))
	}

	//To allocate slice for request body
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		log.Printf("Content-Length failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Read body data to parse json
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Printf("read failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Content-Length:%v", length)

	var geos []ConsumerGeoInfo
	err = json.Unmarshal(body[:length], &geos)
	if err != nil {
		log.Printf("json.Unmarshal failed!! %v", err)
		log.Printf("json:%s", body[:length])
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := connect()
	if err != nil {
		log.Printf("database connect failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("connected database!!")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("transaction begin failed!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	log.Printf("transaction begin!!")

	for _, geo := range geos {
		log.Printf("%v\n", geo)
		if err := addConsumerGeo(db, geo); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("commit faled!!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("commit finished!!")

	w.WriteHeader(http.StatusOK)
}

func Router() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/employees/{1}", employeeHandler)
	//r.HandleFunc("/", Sleeper)
	r.HandleFunc("/consumer/@{latitude:[0-9]+.?[0-9]+},{longtitude:[0-9]+.?[0-9]+}", ConsumerHandler).Methods("GET")
	r.HandleFunc("/consumer/GeoCollection", ConsumerGeoCollectionWriter).Methods("POST")
	return r
}

func main() {
	port := flag.Int("port", 80, "port number")
	dbserver := flag.String("dbserver", db_host_default, "database server")
	dbport := flag.Int("dbport", db_port_default, "database port")
	flag.Parse()
	db_host = *dbserver
	db_port = *dbport

	http.Handle("/", Router())

	log.Printf("Start Go HTTP Server")

	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
