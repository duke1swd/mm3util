package main

import (
	"flag"
	"fmt"
	"log"
)

/*
 * The members object has at least these fields
 */

type member struct {
	address          string
	deliveryMode     string
	displayName      string
	email            string
	listId           string
	memberId         string
	role             string
	selfLink         string
	subscriptionMode string
	user             string
}

var members []member

func membersCmd() {

	if flag.NArg() != 1 {
		fmt.Printf("members command has no subcommands or arguments\n")
		usage()
	}

	setDomain()
	loadMembers()

	for _, v := range members {
		memberDisplay(v)
	}
}

func memberDisplay(m member) {
	fmt.Printf("address: %s\n", m.address)
	fmt.Printf("\tdelivery_mode:     %s\n", m.deliveryMode)
	fmt.Printf("\tdisplay_name:      %s\n", m.displayName)
	fmt.Printf("\temail:             %s\n", m.email)
	fmt.Printf("\tlist_id:           %s\n", m.listId)
	fmt.Printf("\tmember_id:         %s\n", m.memberId)
	fmt.Printf("\trole:              %s\n", m.role)
	fmt.Printf("\tself_link:         %s\n", m.selfLink)
	fmt.Printf("\tsubscription_mode: %s\n", m.subscriptionMode)
	fmt.Printf("\tuser:              %s\n\n", m.user)
}

func loadMembers() {
	res, ok := get(configuration.Url + "/members")
	if !ok {
		log.Fatal("Cannot load members collection")
	}

	resMap, ok := res.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get members")
	}

	entriesRaw, ok := resMap["entries"]
	if !ok {
		log.Fatal("missing entries in members response")
	}

	entries, ok := entriesRaw.([]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get members/entries")
	}

	members = make([]member, len(entries))

	for i, e := range entries {
		entry, ok := e.(map[string]interface{})
		if !ok {
			log.Fatal("badly formed json in response to get members/entries/entry")
		}
		members[i].address = jsonDecode(entry, "members", "address")
		members[i].deliveryMode = jsonDecode(entry, "members", "delivery_mode")
		members[i].displayName = jsonDecode(entry, "members", "display_name")
		members[i].email = jsonDecode(entry, "members", "email")
		members[i].listId = jsonDecode(entry, "members", "list_id")
		members[i].memberId = jsonDecode(entry, "members", "member_id")
		members[i].role = jsonDecode(entry, "members", "role")
		members[i].selfLink = jsonDecode(entry, "members", "self_link")
		members[i].subscriptionMode = jsonDecode(entry, "members", "subscription_mode")
		members[i].user = jsonDecode(entry, "members", "user")
	}
}

func jsonDecode(entry map[string]interface{}, collection string, field string) string {
	raw, ok := entry[field]
	if !ok {
		log.Fatalf("no field \"%s\" in in collection %s", field, collection)
	}
	r, ok := raw.(string)
	if !ok {
		log.Fatalf("badly formed json in response to get %s / %s", collection, field)
	}
	return r
}
