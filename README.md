# dugble-go

Official Go SDK for the Dugble API.

> This SDK is under active development. The current bootstrap branch establishes the shared client and begins the public API surface with Topics and Segments.

## Install

```bash
go get github.com/dugble/dugble-go
```

## Quickstart

```go
package main

import (
	"context"
	"log"

	dugble "github.com/dugble/dugble-go"
)

func main() {
	client := dugble.New("dgb_team_...")

	topics, err := client.Topics.List(context.Background(), dugble.ListTopicsParams{
		Limit: 20,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("topics: %d", len(topics.Data))
}
```

The client defaults to `https://api.dugble.com`, sends bearer authentication automatically, unwraps Dugble API envelopes, and returns structured `*dugble.APIError` values for API failures.
