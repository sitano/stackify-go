# stackify-go
An unofficial Go library for the Stackify metrics/logs

## Current Status
The [Stackify API](https://github.com/stackify/stackify-api/blob/master/endpoints/). It only supports the sending of events.

## Usage Example

```go
package main

import (
	"log"

	"github.com/sitano/stackify-go"
)

func main() {
	if resp, err := stackify.NewClient().Send(
		stackify.CreateReportFromMessages([]*stackify.Event{
			stackify.CreateEvent(stackify.Info, "test"),
	})); err != nil {
		log.Fatalln("Error:", err)
	} else {
		log.Printf("Response: %+v\n", resp)
	}
}

```

## Author
(c) 2015 Ivan Prisyazhnyy
