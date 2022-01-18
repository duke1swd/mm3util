package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type configInfo struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var configuration configInfo

func config() {
	loadConfigFile(flagc)
}

func loadConfigFile(fileName string) {
	raw, er := ioutil.ReadFile(fileName)
	defer func() { raw = nil }() // free up the memory used by the raw, unparsed version of the file

	if er != nil {
		log.Fatalln("Cannot open config file", fileName, "for reading")
	}

	er = json.Unmarshal(raw, &configuration)
	if er != nil {
		log.Printf("Syntax error parsing config file %s", fileName)
		log.Fatal(er)
	}
}
