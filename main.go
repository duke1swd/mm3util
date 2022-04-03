/*
 * This program provides a set of commands for manipulating mailman3 lists
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	flagc string // location of the configuration file
	flagd string // domain name
	flagi bool   // interactive mode
	flagL string // name of the log file
	debug bool
)

func init() {
	flag.StringVar(&flagc, "c", "/opt/mailman/mm/mm3util.cfg", "configuration file name")
	flag.StringVar(&flagd, "d", "", "domain name")
	flag.BoolVar(&flagi, "i", false, "set interactive mode -- errors to stderr not log file")
	flag.StringVar(&flagL, "L", "/var/log/mm3util.log", "log file name")
	flag.BoolVar(&debug, "D", false, "enable debugging")
	flag.Parse()

	if !flagi {
		f, err := os.OpenFile(flagL, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("error opening logfile %s: %v", flagL, err)
		} else {

			log.SetOutput(f)
		}
	}
}

func usage() {
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println("\nCommands:")
	fmt.Println("\tdomains")
	fmt.Println("\tlist")
	fmt.Println("\t\temail <list name> (displays all emails in list)")
	fmt.Println("\tlists")
	fmt.Println("\tuser\tadd <email address>")
	fmt.Println("\t\tdelete <email address|user ID>")
	fmt.Println("\t\tshow <email address|'all'>")
	fmt.Println("\taddress\tadd <email address>")
	fmt.Println("\t\tdelete <email address>")
	fmt.Println("\t\tshow <email address|'all'>")
	fmt.Println("\t\tunlink <email address>")
	fmt.Println("\tmembers")
	fmt.Println("\tsubscribe <list name> <email address")
	fmt.Println("\tunsubscribe <list name> <email address")
	os.Exit(1)
}

func main() {
	config()

	if flag.NArg() < 1 {
		usage()
	}

	switch flag.Arg(0) {
	case "domains":
		domainsCmd()
	case "lists":
		listsCmd()
	case "user":
		userCmd()
	case "address":
		addressCmd()
	case "subscribe":
		subscribeCmd()
	case "unsubscribe":
		unSubscribeCmd()
	case "list":
		listCmd()
	case "members":
		membersCmd()
	default:
		fmt.Println("Unknown command")
		usage()
	}
}
