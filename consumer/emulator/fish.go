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

	geo "github.com/tesujiro/smallfish/consumer/geo"
)

const port = 31080
const host = "127.0.0.1"
const clients = 10
const record_interval_sec = 1
const send_interval_sec = 3
const life_sec = 10

type fish struct {
	id   int
	geos []geo.ConsumerGeoInfo
}

func (f *fish) doRequest(ba []byte) error {
	url := fmt.Sprintf("http://%s:%d/consumer/GeoCollection", host, port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(ba))
	if err != nil {
		fmt.Printf("Request post failed.\n")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read response body failed.\n")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Status Not OK :%s\n", resp.StatusCode)
		return fmt.Errorf("Status Not OK : %d", resp.StatusCode)
	}

	fmt.Printf("response.body:%v\n", body)
	return nil
}

func (f *fish) request() {
	ba, err := json.Marshal(f.geos)
	if err != nil {
		fmt.Printf("json.Marshal failed. error:%v\n", err)
		return
	}
	f.geos = []geo.ConsumerGeoInfo{}

	if err := f.doRequest(ba); err != nil {
		fmt.Printf("http request error:%v", err)
	}
}

func (f *fish) record() {
	f.geos = append(f.geos, geo.ConsumerGeoInfo{ConsumerId: f.id, Timestamp: time.Now(), Lat: 123.456, Lng: 456.789})
}

func (f *fish) walk(ctx context.Context) {
	send_tick := time.NewTicker(time.Second * time.Duration(send_interval_sec)).C
	record_tick := time.NewTicker(time.Second * time.Duration(record_interval_sec)).C
	stop := make(chan bool)
	go func() {
	loop:
		for {
			select {
			case <-record_tick:
				f.record()
			case <-send_tick:
				f.request()
			case <-stop:
				f.request()
				break loop
			}
		}
	}()
	time.Sleep(time.Second * time.Duration(life_sec))
	close(stop)
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
		go func(id int) {
			f := &fish{id: id}
			f.walk(ctx)
			wg.Done()
		}(i)
	}
	wg.Wait()

	return 0
}
