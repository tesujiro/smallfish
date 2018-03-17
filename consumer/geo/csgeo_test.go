package csgeo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {

	config, err := NewConfig()
	if err != nil {
		log.Printf("init config failed: %v", err)
	}

	consumer := NewConsumer(config)

	cases := []struct {
		method string
		body   []ConsumerGeoInfo
		url    string
		err    error
		status int
	}{
		//{url: "", err: nil, status: 200},
		//{method: "GET", url: "/", err: nil, status: http.StatusOK},
		{method: "GET", url: "/xxx", err: nil, status: http.StatusNotFound},
		//{method: "GET", url: "/consumer/@123.456,456.123", body: []ConsumerGeoInfo{}, err: nil, status: http.StatusOK},
		{method: "POST", url: "/consumer/GeoCollection", body: []ConsumerGeoInfo{
			ConsumerGeoInfo{ConsumerId: 1, Timestamp: time.Now(), Lat: 123.456, Lng: 456.789},
			ConsumerGeoInfo{ConsumerId: 1, Timestamp: time.Now().Add(1 * time.Second), Lat: 123.999, Lng: 456.999},
		}, err: nil, status: http.StatusOK},
	}

	for _, c := range cases {
		input, err := json.Marshal(c.body)
		if err != nil {
			t.Fatal(err)
		}

		r, err := http.NewRequest(c.method, c.url, bytes.NewBuffer(input))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Content-Length", strconv.Itoa(len(input)))

		w := httptest.NewRecorder()
		consumer.Router().ServeHTTP(w, r)

		if status := w.Code; status != c.status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, c.status)
		}
		//fmt.Printf("w.Body=%v\n", w.Body.String())
		fmt.Printf("w=%v\n", w)

	}
}
