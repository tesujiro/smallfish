package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Sleeper(w http.ResponseWriter, r *http.Request) {
	log.Printf("Sleeper!!")
	q := r.URL.Query()
	i, _ := strconv.Atoi(q.Get("timer"))
	time.Sleep(time.Duration(i) * time.Millisecond)
	fmt.Fprintf(w, "Hello, World :slept %d msec\n", i)
}

type ConsumerGeoInfo struct {
	ConsumerId int     `json:"consumerId"`
	Lat        float64 `json:"latitude"`
	Lng        float64 `json:"longtitude"`
}

func ConsumerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConsumerHandler!!")
	vars := mux.Vars(r)
	//lat, lng := vars["latitude"], vars["longtitude"]
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

	//fmt.Fprintf(w, "longtitude=%v\n", lng)
	log.Printf("geo=%v\n", geo)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(geo); err != nil {
		log.Print("json.NewEncoder error!\n")
		log.Fatal(err)
	}
	//w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusNotFound)

	//w.Write([]byte("aaabbbccc"))
}

func Router() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/employees/{1}", employeeHandler)
	//r.HandleFunc("/", Sleeper)
	r.HandleFunc("/consumer/@{latitude:[0-9]+.?[0-9]+},{longtitude:[0-9]+.?[0-9]+}", ConsumerHandler)
	return r
}

func main() {
	port := flag.Int("port", 80, "port number")
	flag.Parse()

	http.Handle("/", Router())

	log.Printf("Start Go HTTP Server")

	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
