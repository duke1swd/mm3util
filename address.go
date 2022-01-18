package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type addressT struct {
	Email          string `json:"email"`
	Original_email string `json:"original_email"`
}

func addressCmd() {
	setDomain()

	if flag.NArg() < 2 {
		usage()
	}

	switch flag.Arg(1) {
	case "add":
		addressAdd()
	case "delete":
		addressDelete()
	case "show":
		addressShow()
	case "unlink":
		addressUnlink()
	default:
		fmt.Printf("Unknown address subcommand %s\n", flag.Arg(1))
		usage()
	}

}

func addressAdd() {
	if flag.NArg() != 3 {
		fmt.Println("address add requires one address")
		usage()
	}

	address := flag.Arg(2)
	resRaw, ok := get(configuration.Url + "/addresses/" + address)
	if ok {
		res, ok := resRaw.(map[string]interface{})
		if !ok {
			log.Fatal("badly formed json in response to get addresses/address (add)")
		}

		_, ok = res["self_link"]
		if ok {
			log.Printf("Address %s already exists, not added.", address)
			os.Exit(2) // error code 2 == benign failure
		}
	}

	var addressStruct addressT
	addressStruct.Email = address
	addressStruct.Original_email = address
	post(configuration.Url+"/addresses/"+address, addressStruct)
	log.Printf("Address %s added", address)
}

func addressDelete() {
	if flag.NArg() != 3 {
		fmt.Println("address delete requires one address")
		usage()
	}

	address := flag.Arg(2)

	resRaw, ok := get(configuration.Url + "/addresses/" + address)
	if !ok {
		log.Printf("Address %s not found, not deleted", address)
		os.Exit(2)
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get addresses/address (del)")
	}

	linkRaw, ok := res["self_link"]
	if !ok {
		log.Printf("Address %s not found, not deleted", address)
		os.Exit(2)
	}

	link, ok := linkRaw.(string)
	if !ok {
		log.Fatalf("Address %s has a non-string self link (del)\n", address)
	}

	delCmd(link)
	log.Printf("Address %s deleted.", address)
}

func addressShow() {
	if flag.NArg() < 3 {
		fmt.Println("address show requires either an address or \"all\"")
		usage()
	}

	address := flag.Arg(2)
	if address == "all" || address == "All" || address == "ALL" {
		addressShowAll()
	} else {
		addressShowOne(address)
	}
}

func addressShowAll() {
	resRaw, ok := get(configuration.Url + "/addresses")
	if !ok {
		log.Fatal("unable to get addresses collection")
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get addresses")
	}

	entriesRaw, ok := res["entries"]
	if !ok {
		log.Fatal("no entries in response to get addresses")
	}
	entries, ok := entriesRaw.([]interface{})
	if !ok {
		log.Fatal("badly formed json /entries in response to get addresses")
	}

	for _, entryRaw := range entries {
		entry, ok := entryRaw.(map[string]interface{})
		if !ok {
			log.Fatalf("badly formed address entry: %v\n", entryRaw)
		}

		addressDisplay(entry)
	}
}

func addressDisplay(addr map[string]interface{}) {
	for key, value := range addr {
		if key != "http_etag" && key != "self_link" {
			fmt.Printf("\t%s: %v\n", key, value)
		}
	}
	fmt.Println("")
}

func addressShowOne(address string) {
	resRaw, ok := get(configuration.Url + "/addresses/" + address)
	if !ok {
		log.Fatalf("address %s not found", address)
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get addresses/address")
	}

	addressDisplay(res)
}

func addressUnlink() {
	if flag.NArg() != 3 {
		fmt.Println("address unlink requires one address")
		usage()
	}

	address := flag.Arg(2)
	resRaw, ok := get(configuration.Url + "/addresses/" + address)
	if !ok {
		log.Fatalf("address %s not found.  Not unlinked", address)
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get addresses/address (unlink)")
	}

	_, ok = res["self_link"]
	if !ok {
		log.Fatalf("address %s not found.  Not unlinked (self link)", address)
	}

	delCmd(configuration.Url + "/addresses/" + address + "/user")
	log.Printf("address %s unlinked", address)
}
