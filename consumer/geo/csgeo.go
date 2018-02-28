package main

import (
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

func ConsumerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConsumerHandler!!")
	vars := mux.Vars(r)

	fmt.Fprintf(w, "URL=%v\n", r.URL)
	fmt.Fprintf(w, "Vars=%v\n", vars)
	fmt.Fprintf(w, "longtitude=%v\n", vars["longtitude"])
	fmt.Fprintf(w, "latitude=%v\n", vars["latitude"])
	//log.Printf("URL=%v\n", r.URL)
	log.Printf("Vars=%v\n", vars)
	log.Printf("longtitude=%v\n", vars["longtitude"])
	log.Printf("latitude=%v\n", vars["latitude"])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func Router() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/employees/{1}", employeeHandler)
	r.HandleFunc("/consumer/@{longtitude:[0-9]+},{latitude:[0-9]+}", ConsumerHandler)
	return r
}

func main() {
	port := flag.Int("port", 80, "port number")
	flag.Parse()

	//r := mux.NewRouter()
	//r.HandleFunc("/", Sleeper)
	//r.HandleFunc("/consumer/@{longtitude:[0-9]+},{latitude:[0-9]+}", ConsumerHandler)
	//http.Handle("/", r)
	http.Handle("/", Router())

	log.Printf("Start Go HTTP Server")

	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
