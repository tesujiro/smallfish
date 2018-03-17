package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler(t *testing.T) {

	cases := []struct {
		method string
		body   []ConsumerGeoInfo
		url    string
		err    error
		status int
	}{
		//{url: "", err: nil, status: 200},
		//{method: "GET", url: "/", err: nil, status: http.StatusOK},
		//{method: "GET", url: "/xxx", err: nil, status: http.StatusNotFound},
		{method: "GET", url: "/consumer/@123.456,456.123", body: []ConsumerGeoInfo{}, err: nil, status: http.StatusOK},
		//{method: "POST", url: "/consumer/GeoCollection", body: `{geo:[{lat:111,lng:222},{lat:333,lng:444}]}`, err: nil, status: http.StatusOK},
		//{method: "POST", url: "/consumer/GeoCollection", body: `[{consumerId:1,latitude:12.34,longtitud:45.67}]`, err: nil, status: http.StatusOK},
		{method: "POST", url: "/consumer/GeoCollection", body: []ConsumerGeoInfo{
			ConsumerGeoInfo{ConsumerId: 1, Lat: 123.456, Lng: 456.789},
			ConsumerGeoInfo{ConsumerId: 1, Lat: 123.999, Lng: 456.999},
		}, err: nil, status: http.StatusOK},
		//{method: "GET", url: "/consumer/@123.11,456.23", body: ,err: nil, status: http.StatusOK},
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
		Router().ServeHTTP(w, r)

		if status := w.Code; status != c.status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, c.status)
		}
		//fmt.Printf("w.Body=%v\n", w.Body.String())
		fmt.Printf("w=%v\n", w)

	}
}
