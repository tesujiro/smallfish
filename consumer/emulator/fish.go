package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	//"github.com/tesujiro/smallfish/consumer/geo"
	geo "github.com/tesujiro/smallfish/consumer/geo"
)

const port = 31080
const host = "127.0.0.1"
const clients = 10
const request_interval_sec = 1

func doRequest(ba []byte) error {
	url := fmt.Sprintf("http://%s:%d/consumer/GeoCollection", host, port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(ba))
	if err != nil {
		fmt.Printf("Request post failed.")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read response body failed.")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Status Not OK :%s", resp.StatusCode)
		return fmt.Errorf("Status Not OK : %d", resp.StatusCode)
	}

	fmt.Printf("response.body:%v\n", body)
	return nil
}

func request(id int, ctx context.Context) {
	geos := []geo.ConsumerGeoInfo{
		geo.ConsumerGeoInfo{ConsumerId: id, Timestamp: time.Now(), Lat: 123.456, Lng: 456.789},
	}
	ba, err := json.Marshal(geos)
	if err != nil {
		fmt.Printf("json.Marshal failed. error:%v", err)
		return
	}

	if err := doRequest(ba); err != nil {
		fmt.Printf("http request error:%v", err)
	}
}

func main() {
	//csgeo.Run()
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error:\n%s", err)
			os.Exit(1)
		}
	}()
	os.Exit(_main())
}

func _main() int {
	if envvar := os.Getenv("GOMAXPROCS"); envvar == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for i := 1; i <= clients; i++ {
		wg.Add(1)
		id := i
		go func() {
			request(id, ctx)
			wg.Done()
		}()
	}
	wg.Wait()

	return 0
}
