package main

import (
	"flag"
	"fmt"
	"log"
)

/*
	The domains object has 4 top level objects of interest
	"entries" which is all the stuff
	"start" which should always be zero, and
	"total_size" which is an integer, the number of domains
	"http_etag" which is what?

	Entries is an array of maps.  Entry maps have these keys
		alias_domain	not sure what this is
		description	some text
		mail_host	this is the domain's name
		self_link	the url to the domain's REST api
		http_etag1

	The info in the entries is copied into the domains array.
*/

type domain struct {
	mailHost    string
	description string
	url         string
}

var (
	domainId int
	domains  []domain
)

func domainsCmd() {
	loadDomains()

	if flag.NArg() != 1 {
		fmt.Printf("domains command has no subcommands or arguments\n")
		usage()
	}

	for _, v := range domains {
		fmt.Printf("%s %s\n", v.mailHost, v.description)
	}
}

func setDomain() {
	loadDomains()

	domainId = -1

	if flagd == "" {
		domainId = 0
	} else {
		for i, v := range domains {
			if v.mailHost == flagd {
				domainId = i
				break
			}
		}
		if domainId < 0 {
			log.Fatalf("cannot find domain %s", flagd)
		}
	}
}

func loadDomains() {
	res, ok := get(configuration.Url + "/domains")
	if !ok {
		log.Fatal("Cannot load domains collection")
	}

	resMap, ok := res.(map[string]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get domains")
	}

	entriesRaw, ok := resMap["entries"]
	if !ok {
		log.Fatal("missing entries in domain response")
	}

	entries, ok := entriesRaw.([]interface{})
	if !ok {
		log.Fatal("badly formed json in response to get domains/entries")
	}

	domains = make([]domain, len(entries))

	for i, e := range entries {
		//fmt.Printf("Entry %d is %v\n\n", k, e)
		entry, ok := e.(map[string]interface{})
		if !ok {
			log.Fatal("badly formed json in response to get domains/entries/entry")
		}
		mhRaw, ok := entry["mail_host"]
		if !ok {
			log.Fatal("missing mail_host")
		}
		domains[i].mailHost, ok = mhRaw.(string)
		if !ok {
			log.Fatal("mail_host not a string")
		}

		desRaw, ok := entry["description"]
		if !ok {
			log.Fatal("missing description")
		}
		domains[i].description, ok = desRaw.(string)
		if !ok {
			log.Fatal("description not a string")
		}

		urlRaw, ok := entry["self_link"]
		if !ok {
			log.Fatal("missing self_link")
		}
		domains[i].url, ok = urlRaw.(string)
		if !ok {
			log.Fatal("self_link not a string")
		}
	}
}
