package csgeo

import (
	"context"
	"database/sql"
	"encoding/json"
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

type Consumer struct {
	config Config
}

func NewConsumer(c *Config) *Consumer {
	return &Consumer{config: *c}
}

type ConsumerGeoInfo struct {
	ConsumerId int       `json:"consumerId"`
	Timestamp  time.Time `json:"timestamp"`
	Lat        float64   `json:"latitude"`
	Lng        float64   `json:"longtitude"`
}

func (c *Consumer) connect() (*sql.DB, error) {
	url := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable", c.config.db_user, c.config.db_host, c.config.db_port, c.config.db_name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		//log.Fatal("error connecting to the database: ", err)
		fmt.Println("error connecting to the database: ", err)
	}
	return db, err
}

func (c *Consumer) addConsumerGeo(db *sql.DB, geo ConsumerGeoInfo) error {
	// Insert two rows into the "location" table.
	stmt, err := db.Prepare("INSERT INTO location (id, time, lat, lng) VALUES ($1,$2,$3,$4)")
	if err != nil {
		log.Printf("prepare statement faled!! %v", err)
		return err
	}
	defer stmt.Close() // danger!
	log.Printf("prepare statement!!")

	res, err := stmt.Exec(geo.ConsumerId, geo.Timestamp, geo.Lat, geo.Lng)
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

/*
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

	// CHU RYAKU!!

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusNotFound)

	if err := json.NewEncoder(w).Encode(geo); err != nil {
		log.Print("json.NewEncoder error!\n")
		log.Fatal(err)
	}
}
*/

func (c *Consumer) GeoCollectionWriter(w http.ResponseWriter, r *http.Request) {
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

	db, err := c.connect()
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
		if err := c.addConsumerGeo(db, geo); err != nil {
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

func (c *Consumer) Router() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/", Sleeper)
	//r.HandleFunc("/consumer/@{latitude:[0-9]+.?[0-9]+},{longtitude:[0-9]+.?[0-9]+}", ConsumerHandler).Methods("GET")
	r.HandleFunc("/consumer/GeoCollection", c.GeoCollectionWriter).Methods("POST")
	return r
}

func Run(ctx context.Context) {

	config, err := NewConfig()
	if err != nil {
		log.Printf("init config failed: %v", err)
	}

	consumer := NewConsumer(config)

	http.Handle("/", consumer.Router())

	log.Printf("Start Go HTTP Server")

	err = http.ListenAndServe(":"+strconv.Itoa(consumer.config.http_port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
