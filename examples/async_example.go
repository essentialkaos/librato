package main

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/essentialkaos/librato"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	librato.Mail = "mail@domain.com"
	librato.Token = "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234"

	// Create struct for async sending metrics data
	// With this preferences metrics will be sended to Librato once
	// a minute or if queue size reached 60 elements
	metrics, err := librato.NewMetrics(time.Minute, 60)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for {
		metrics.Add(
			librato.Gauge{
				Name:  "example:gauge_1",
				Value: randomInt(1000),
			},
		)

		time.Sleep(15 * time.Second)
	}
}

func randomInt(n int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(n)
}
