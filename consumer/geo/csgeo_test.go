package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {

	cases := []struct {
		method string
		url    string
		err    error
		status int
	}{
		//{url: "", err: nil, status: 200},
		//{method: "GET", url: "/", err: nil, status: http.StatusOK},
		//{method: "GET", url: "/xxx", err: nil, status: http.StatusNotFound},
		{method: "GET", url: "/consumer/@123,456", err: nil, status: http.StatusOK},
		//{method: "GET", url: "/consumer/@123.11,456.23", err: nil, status: http.StatusOK},
	}

	for _, c := range cases {
		r, err := http.NewRequest(c.method, c.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		Router().ServeHTTP(w, r)

		if status := w.Code; status != c.status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, c.status)
		}
		//fmt.Printf("w.Body=%v\n", w.Body.String())
		fmt.Printf("w=%v\n", w)

	}
}
