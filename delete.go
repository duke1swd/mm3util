package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func delCmd(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("http.NewRequest(\"DELETE\"... fails with %s", err.Error())
	}

	// set up the headers
	req.SetBasicAuth(configuration.Username, configuration.Password)

	// Maket the API call
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("iclient.Do (DEL) fails with %s", err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll (DEL) fails with %s", err.Error())
	}

	if debug {
		fmt.Printf("http DEL request \"%s\" returns %d\n", url, resp.StatusCode)
		prettyJ("JSON from delete:", bodyBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("delete of \"%s\" fails with code %d", url, resp.StatusCode)
	}
}
