package main

import (
	"flag"
	"log"

	"github.com/sitano/stackify-go"
)

func main() {
	flag.Parse()

	if resp, err := stackify.NewClient().Send(
		stackify.CreateReportFromMessages([]*stackify.Event{
			stackify.CreateEvent(stackify.Info, "test"),
	})); err != nil {
		log.Fatalln("Error:", err)
	} else {
		log.Printf("Response: %+v\n", resp)
	}
}
