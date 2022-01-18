package main

import (
	"flag"
	"fmt"
	"log"
)

type subscribeT struct {
	ListID         string `json:"list_id"`
	Subscriber     string `json:"subscriber"`
	PreVerified    bool   `json:"pre_verified"`
	PreConfirmed   bool   `json:"pre_confirmed"`
	PreAproved     bool   `json:"pre_Approved"`
	SendWelcome    bool   `json:"send_welcome_message"`
	DeliveryMode   string `json:"delivery_mode"`
	DeliveryStatus string `json:"delivery_status"`
}

func subscribeCmd() {
	setDomain()

	if flag.NArg() < 3 {
		fmt.Println("subscribe requires a list and an email address")
		usage()
	}

	list := flag.Arg(1)
	address := flag.Arg(2)

	if flag.NArg() > 3 {
		fmt.Println("too many arguments to user add -- options NYI")
		usage()
	}

	listEntry := findList(list)

	resRaw, ok := get(configuration.Url + "/addresses/" + address)
	if !ok {
		log.Fatalf("Address %s not found, not subsdribed", address)
	}

	res, ok := resRaw.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get addresses/address (subscribe)")
	}

	_, ok = res["self_link"]
	if !ok {
		log.Fatalf("Address %s not found, not subsdribed (self-link)", address)
		return
	}

	var subscribeStruct subscribeT
	subscribeStruct.ListID = listEntry.listId
	subscribeStruct.Subscriber = address
	subscribeStruct.PreVerified = true
	subscribeStruct.PreConfirmed = true
	subscribeStruct.PreAproved = true
	subscribeStruct.SendWelcome = false
	subscribeStruct.DeliveryMode = "regular"
	subscribeStruct.DeliveryStatus = "by_user"

	post(configuration.Url+"/members", subscribeStruct)
	log.Printf("user %s subscribed to list %s", address, list)
}
