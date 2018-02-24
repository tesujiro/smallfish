package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// http.HandleFuncに登録する関数
// http.ResponseWriterとhttp.Requestを受ける
func Sleeper(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	i, _ := strconv.Atoi(q.Get("timer"))
	time.Sleep(time.Duration(i) * time.Millisecond)
	fmt.Fprintf(w, "Hello, World :slept %d msec\n", i)
}

func main() {
	port := flag.Int("port", 80, "port number")
	flag.Parse()

	http.HandleFunc("/", Sleeper)

	log.Printf("Start Go HTTP Server")

	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
