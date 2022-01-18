package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func prettyJ(title string, bodyBytes []byte) {
	fmt.Println(title)

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, bodyBytes, "", "    "); err != nil {
		fmt.Printf("Bad JSON: err = %v\n", err)
	} else {
		fmt.Println(prettyJSON.String())
	}
}
