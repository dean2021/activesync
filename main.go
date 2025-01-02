package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dean2021/activesync/ashttp"
)

func main() {
	query, err := ashttp.ParseBase64Query("oQkECBCeDEK6NjuTWKLjgUH2WCxdBIIanKgLV2luZG93c01haWw=")
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.MarshalIndent(query, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
