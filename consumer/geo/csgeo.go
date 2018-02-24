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
	q := r.URL.Query()
	i, _ := strconv.Atoi(q.Get("timer"))
	time.Sleep(time.Duration(i) * time.Millisecond)
	fmt.Fprintf(w, "Hello, World :slept %d msec\n", i)
}

func ConsumerHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ConsumerHandler")
	vars := mux.Vars(r)

	fmt.Fprintf(w, "URL=%v\n", r.URL)
	fmt.Fprintf(w, "Vars=%v\n", vars)
}

func main() {
	port := flag.Int("port", 80, "port number")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", Sleeper)
	r.HandleFunc("/consumer/@{longtitude:[0-9]+},{latitude:[0-9]+}", ConsumerHandler)
	http.Handle("/", r)

	log.Printf("Start Go HTTP Server")

	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
