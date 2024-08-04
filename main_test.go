package main

import (
	"fmt"
	"testing"
	"github.com/gorilla/mux"
)


func Test_readUrl(*T testing.T) {
	r := mux.NewRouter()
	r.Handle("/latest/{currency}", readUrl())
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/latest/USD")
	if err != nil {
		t.Errorf("Error: %v ", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Error: %v ", resp.Status)
	}

	// Process the request for the version6 API
	// Unmarshall the response body into the ApiResponse struct
	ApiResponse := &ApiV6Response{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v ", err)
	}
}
