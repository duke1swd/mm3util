package main

import (
	"flag"
	"fmt"
	"log"
)

/*
	The list object has these
*/

type listT struct {
	displayName  string
	listName     string
	mailHost     string
	advertised   bool // not yet implemented
	listId       string
	memberCount  int64 // not yet implemented
	volume       int64 // not yet implemented
	description  string
	selfLink     string
	fqdnListname string
}

var lists []listT

func listsCmd() {
	setDomain()
	loadLists()

	if flag.NArg() != 1 {
		fmt.Printf("lists command has no subcommands or arguments\n")
		usage()
	}

	for _, v := range lists {
		fmt.Printf("%s \"%s\"\n", v.displayName, v.description)
	}
}

func findList(list string) listT {
	loadLists()
	for _, l := range lists {
		if l.listName == list {
			return l
		}
	}
	log.Fatalf("list %s not found", list)
	panic("not reached")
}

func loadLists() {
	res, ok := get(domains[domainId].url + "/lists")
	if !ok {
		log.Fatal("Cannot load lists collection")
	}

	resMap, ok := res.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get lists")
	}

	entriesRaw, ok := resMap["entries"]
	if !ok {
		log.Fatal("missing entries in list response")
	}

	entries, ok := entriesRaw.([]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get lists/entries")
	}

	lists = make([]listT, len(entries))

	for i, e := range entries {
		entry, ok := e.(map[string]interface{})
		if !ok {
			log.Fatal("badly formed json in response to get lists/entries/entry")
		}

		lists[i].displayName = jsonDecode(entry, "lists", "display_name")
		lists[i].listName = jsonDecode(entry, "lists", "list_name")
		lists[i].mailHost = jsonDecode(entry, "lists", "mail_host")
		lists[i].description = jsonDecode(entry, "lists", "description")
		lists[i].listId = jsonDecode(entry, "lists", "list_id")
		lists[i].selfLink = jsonDecode(entry, "lists", "self_link")
		lists[i].fqdnListname = jsonDecode(entry, "lists", "fqdn_listname")
	}
}
