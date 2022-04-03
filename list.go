package main

import (
	"flag"
	"fmt"
	"log"
)

/*
	The roster of an email list has these, plus a few more
*/

type rosterT struct {
	address           string
	bounceScore       int64 // not yet implemented
	deliveryMode      string
	displayName       string
	email             string
	listId            string
	memberId          int64 // not yet implemented
	role              string
	selfLink          string
	subscriptionMode  string
	totalWarningsSent int64 // not yet implemented
	user              string
}

var roster []rosterT

func listCmd() {
	setDomain()
	loadLists()

	if flag.NArg() != 3 {
		usage()
	}

	switch flag.Arg(1) {
	case "email":
		listEmail()
	case "emails":
		listEmail()
	case "name":
		listName()
	case "names":
		listName()
	default:
		fmt.Printf("Unknown list subcommand %s\n", flag.Arg(1))
		usage()
	}
}

func listEmail() {
	listLoad()

	for _, r := range roster {
		fmt.Println(r.email)
	}
}

func listName() {
	listLoad()

	for _, r := range roster {
		fmt.Println(r.displayName)
	}
}

func listLoad() {
	listName := flag.Arg(2)
	var list listT
	list.listName = ""

	for _, i := range lists {
		if listName == i.displayName ||
			listName == i.listName ||
			listName == i.listId ||
			listName == i.fqdnListname {
			list = i
			break
		}
	}

	if list.listName == "" {
		fmt.Printf("List %s not found\n", listName)
		usage()
	}

	loadEntries(list.listId)
}

func loadEntries(listId string) {
	res, ok := get(configuration.Url + "/lists/" + listId + "/roster/member")
	if !ok {
		log.Fatalf("Cannot load members of list %s", listId)
	}

	resMap, ok := res.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get lists roster member")
	}

	entriesRaw, ok := resMap["entries"]
	if !ok {
		log.Fatal("missing entries in list roster response")
	}

	entries, ok := entriesRaw.([]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get list roster")
	}

	roster = make([]rosterT, len(entries))

	for i, e := range entries {
		entry, ok := e.(map[string]interface{})
		if !ok {
			log.Fatal("badly formed json in response to get list roster/entries/entry")
		}
		roster[i].address = jsonDecode(entry, "list roster", "address")
		roster[i].deliveryMode = jsonDecode(entry, "list roster", "delivery_mode")
		roster[i].displayName = jsonDecode(entry, "list roster", "display_name")
		roster[i].email = jsonDecode(entry, "list roster", "email")
		roster[i].listId = jsonDecode(entry, "list roster", "list_id")
		roster[i].role = jsonDecode(entry, "list roster", "role")
		roster[i].selfLink = jsonDecode(entry, "list roster", "self_link")
		roster[i].subscriptionMode = jsonDecode(entry, "list roster", "subscription_mode")
		roster[i].user = jsonDecode(entry, "list roster", "user")
	}
}
