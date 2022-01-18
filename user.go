package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func userCmd() {
	setDomain()

	if flag.NArg() < 2 {
		usage()
	}

	switch flag.Arg(1) {
	case "add":
		userAdd()
	case "delete":
		userDelete()
	case "show":
		userShow()
	default:
		fmt.Printf("Unknown user subcommand %s\n", flag.Arg(1))
		usage()
	}

}

type userDT struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type userT struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

func userAdd() {
	if flag.NArg() < 3 {
		fmt.Println("user add requires one address")
		usage()
	}

	address := flag.Arg(2)

	if flag.NArg() > 4 {
		fmt.Println("too many arguments to user add")
		usage()
	}

	var displayName string

	displayNamePresent := false
	if flag.NArg() == 4 {
		displayName = flag.Arg(3)
		displayNamePresent = true
	}

	resRaw, ok := get(configuration.Url + "/addresses/" + address)
	if ok {
		res, ok := resRaw.(map[string]interface{})
		if !ok {
			log.Fatal("badly formed json in response to get addresses/address (add)")
		}

		_, ok = res["self_link"]
		if ok {
			log.Fatalf("Address %s already exists.  User not added\n", address)
		}
	}

	if displayNamePresent {
		var userStruct userDT
		userStruct.Email = address
		userStruct.DisplayName = displayName

		post(configuration.Url+"/users", userStruct)
	} else {
		var userStruct userT
		userStruct.Email = address

		post(configuration.Url+"/users", userStruct)
	}
}

func userDelete() {
	if flag.NArg() < 3 {
		fmt.Println("user delete requires one address or userID")
		usage()
	}

	address := flag.Arg(2)

	if flag.NArg() > 3 {
		fmt.Println("too many arguments to user delete")
		usage()
	}

	resRaw, ok := get(configuration.Url + "/users/" + address)
	if !ok {
		log.Printf("User %s does not exist\n", address)
		os.Exit(2)
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get users/address (delete)")
	}

	_, ok = res["self_link"]
	if !ok {
		log.Printf("User %s does not exist (self-link)\n", address)
		os.Exit(2)
	}
	delcmd(configuration.Url + "/users/" + address)
	log.Printf("User %s deleted", address)
}

func userShow() {
	if flag.NArg() < 3 {
		fmt.Println("user show requires either a user's address or \"all\"")
		usage()
	}

	user := flag.Arg(2)
	if user == "all" || user == "All" || user == "ALL" {
		userShowAll()
	} else {
		userShowOne(user)
	}
}

func userShowAll() {
	resRaw, ok := get(configuration.Url + "/users")
	if !ok {
		log.Fatal("Cannot get users collection")
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get users")
	}

	entriesRaw, ok := res["entries"]
	if !ok {
		log.Fatal("no entries in response to get users")
	}
	entries, ok := entriesRaw.([]interface{})
	if !ok {
		log.Fatal("badly formed json /entries in response to get users")
	}

	for _, entryRaw := range entries {
		entry, ok := entryRaw.(map[string]interface{})
		if !ok {
			log.Fatalf("badly formed user entry: %v\n", entryRaw)
		}

		userDisplay(entry)
	}
}

func userShowOne(user string) {
	resRaw, ok := get(configuration.Url + "/users/" + user)
	if !ok {
		log.Fatalf("user %s not found", user)
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get addresses/address")
	}

	userDisplay(res)
}

func userDisplay(entry map[string]interface{}) {
	for key, value := range entry {
		if key == "password" {
			fmt.Printf("\t%s: <elided>\n", key)
		} else if key != "http_etag" && key != "self_link" {
			fmt.Printf("\t%s: %v\n", key, value)
		}
	}
	fmt.Println("")
}
