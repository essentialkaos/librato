package main

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"

	"pkg.re/essentialkaos/ek.v6/rand"

	"pkg.re/essentialkaos/librato.v4"
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
				Value: rand.Int(1000),
			},
		)

		time.Sleep(15 * time.Second)
	}
}
