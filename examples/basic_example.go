package main

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"

	"pkg.re/essentialkaos/ek.v8/rand"

	"pkg.re/essentialkaos/librato.v6"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	librato.Mail = "mail@domain.com"
	librato.Token = "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234"

	for {
		errs := librato.AddMetric(
			librato.Gauge{
				Name:  "example:gauge_1",
				Value: rand.Int(1000),
			},
			librato.Gauge{
				Name:   "example:gauge_2",
				Value:  float64(rand.Int(1000)) / 5.0,
				Source: "go_librato_example",
			},
			librato.Counter{
				Name:  "example:counter_1",
				Value: rand.Int(1000),
			},
		)

		if len(errs) != 0 {
			fmt.Println("Errors:")

			for _, err := range errs {
				fmt.Printf("  %v\n", err)
			}
		} else {
			fmt.Println("Data sended to Librato Metrics")
		}

		time.Sleep(time.Minute)
	}
}
