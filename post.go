package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func post(url string, payload interface{}) {
	jsonReq, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	if debug {
		fmt.Printf("Posting to url %s\n", url)
		prettyJ("POST payload:", jsonReq)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))

	// set up the headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(configuration.Username, configuration.Password)

	// Make the API call
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("put of \"%s\" fails with code %d", url, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	if debug {
		prettyJ("JSON from POST:", bodyBytes)
	}
}
