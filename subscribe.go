package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type subscribeT struct {
	ListID      string `json:"list_id"`
	Subscriber  string `json:"subscriber"`
	PreAproved  string `json:"pre_approved"`
	PreVerified string `json:"pre_verified"`
	//PreConfirmed string `json:"pre_confirmed"`		// 400
	//DisplayName  string `json:"display_name"`		// 400
	//SendWelcome    string `json:"send_welcome_message"`	// 400
	//DeliveryMode   string `json:"delivery_mode"`		// 400
	//DeliveryStatus string `json:"delivery_status"`	// 400
}

var (
	address string
	list    string
)

func sSetup() (listId, userId string) {
	setDomain()

	if flag.NArg() < 3 {
		fmt.Println("subscribe/unsubscribe requires a list and an email address")
		usage()
	}

	list = flag.Arg(1)
	address = flag.Arg(2)

	if flag.NArg() > 3 {
		fmt.Println("too many arguments to subscribe/unsubscribe")
		usage()
	}

	listEntry := findList(list)
	listId = listEntry.listId

	resRaw, ok := get(configuration.Url + "/users/" + address)
	if !ok {
		log.Fatalf("User %s not found, not subscribed/unsubscribed", address)
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get users/address (subscribe/unsubscribe)")
	}

	_, ok = res["self_link"]
	if !ok {
		log.Fatalf("User %s not found (subscribe/unsubscribe, self-link)", address)
	}

	// Get the user ID
	userIdraw, ok := res["user_id"]
	if !ok {
		log.Fatalf("User %s does not have a user_id (subscribe/unsubscribe)", address)
	}
	userId, ok = userIdraw.(string)
	if !ok {
		log.Fatalf("User %s has bad json (subscribe/unsubscribe)", address)
	}

	return
}

func subscribeCmd() {
	listId, userId := sSetup()

	var subscribeStruct subscribeT
	subscribeStruct.ListID = listId
	subscribeStruct.Subscriber = userId
	subscribeStruct.PreAproved = "true"
	subscribeStruct.PreVerified = "true"
	//subscribeStruct.PreConfirmed = "true"
	//subscribeStruct.DisplayName = "Display Name Foo"
	//subscribeStruct.SendWelcome = "false"
	//subscribeStruct.DeliveryMode = "regular"
	//subscribeStruct.DeliveryStatus = "by_user"

	post(configuration.Url+"/members", subscribeStruct)
	log.Printf("user %s subscribed to list %s", address, list)
}

func unSubscribeCmd() {
	listId, _ := sSetup()

	type memberFinderT struct {
		ListId     string `json:"list_id"`
		Subscriber string `json:"subscriber"`
	}

	var memberFinder memberFinderT

	memberFinder.ListId = listId
	memberFinder.Subscriber = address

	resRaw := post(configuration.Url+"/members/find", memberFinder)
	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get members/find (unsubscribe)")
	}

	totalSizeRaw, ok := res["total_size"]
	if !ok {
		log.Fatal("no total size in response to members/find (unsubscribe)")
	}
	totalSizeF, ok := totalSizeRaw.(float64)
	if !ok {
		log.Fatal("totalSize is not an int64 (unsubscribe)")
	}
	totalSize := int(totalSizeF)

	if totalSize == 0 {
		log.Printf("User %s does not appear to be in list %s.  Not unsubscribed\n",
			address, list)
		os.Exit(2)
	}

	if totalSize != 1 {
		log.Fatalf("User %s appears to be subscribed multiple times in list %s.  Not unsubscribed\n",
			address, list)
	}

	entriesRaw, ok := res["entries"]
	if !ok {
		log.Fatal("no entries in response to members/find")
	}
	entries, ok := entriesRaw.([]interface{})
	if !ok {
		log.Fatal("badly formed json /entries in response to members/find")
	}
	entry, ok := entries[0].(map[string]interface{})
	if !ok {
		log.Fatal("entries[0] bad type in unsubscribe")
	}

	selfLinkRaw, ok := entry["self_link"]
	if !ok {
		log.Fatal("self_link field not found (unsubscribe)")
	}

	selfLink, ok := selfLinkRaw.(string)
	if !ok {
		log.Fatal("self link not string")
	}

	delCmd(selfLink)
	log.Printf("user %s removed from list %s", address, list)
}
