package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	//"github.com/tesujiro/smallfish/consumer/geo"
	g "github.com/tesujiro/smallfish/consumer/geo"
)

const port = 31080
const host = "127.0.0.1"
const fish = 10
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

func main() {
	geos := []g.ConsumerGeoInfo{
		g.ConsumerGeoInfo{ConsumerId: 1, Timestamp: time.Now(), Lat: 123.456, Lng: 456.789},
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
