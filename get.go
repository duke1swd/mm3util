package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Returns unmashalled json, true on success, or nil, false on failure.
func get(url string) (interface{}, bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("http.NewRequest(\"GET\"... fails with %s", err.Error())
	}

	// set up the headers
	//req.Header.Add("Accept", "application/json")
	//req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(configuration.Username, configuration.Password)

	// Maket the API call
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("iclient.Do (GET) fails with %s", err.Error())
	}
	defer resp.Body.Close()

	if debug {
		fmt.Printf("http GET request \"%s\" returns %d\n", url, resp.StatusCode)
	}

	if resp.StatusCode != 200 {
		return nil, false
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll (GET) fails with %s", err.Error())
	}

	if debug {
		prettyJ("JSON from get:", bodyBytes)
	}

	var responseObject interface{}
	json.Unmarshal(bodyBytes, &responseObject)

	return responseObject, true
}
