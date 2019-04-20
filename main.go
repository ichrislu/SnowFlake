package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ichrislu/snowflake/core"
)

var node *snowflake.Node

func main() {
	addrStr := flag.String("addr", ":8080", "TCP network address addr, e.g. 0.0.0.0:8080")
	nodeNum := flag.Int64("node", 0, "node number, between 0 and 1023")

	flag.Parse()

	// Create a new Node with a Node number of 1
	var err error
	node, err = snowflake.NewNode(*nodeNum)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", getId)
	errInfo := http.ListenAndServe(*addrStr, nil)
	if errInfo != nil {
		log.Fatal("ListenAndServe: ", errInfo)
	}
}

func getId(rw http.ResponseWriter, req *http.Request) {
	// node.Generate()
	fmt.Fprint(rw, node.Generate())
	/*
		// Generate a snowflake ID.
		id := node.Generate()

		// Print out the ID in a few different ways.
		fmt.Printf("Int64  ID: %d\n", id)
		fmt.Printf("String ID: %s\n", id)
		fmt.Printf("Base2  ID: %s\n", id.Base2())
		fmt.Printf("Base64 ID: %s\n", id.Base64())

		// Print out the ID's timestamp
		fmt.Printf("ID Time  : %d\n", id.Time())

		// Print out the ID's node number
		fmt.Printf("ID Node  : %d\n", id.Node())

		// Print out the ID's sequence number
		fmt.Printf("ID Step  : %d\n", id.Step())

		// Generate and print, all in one.
		fmt.Printf("ID       : %d\n", node.Generate().Int64())
	*/
}
