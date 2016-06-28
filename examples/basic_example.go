package main

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"

	"pkg.re/essentialkaos/ek.v2/rand"

	"github.com/essentialkaos/librato"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	librato.Mail = "mail@domain.com"
	librato.Token = "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234"

	// We use prefix "example:" which will be added to each gauge and counter sended
	// to Librato
	librato.Prefix = "example:"

	for {
		errs := librato.AddMetric(
			librato.Gauge{
				Name:  "gauge_1",
				Value: rand.Int(1000),
			},
			librato.Gauge{
				Name:   "gauge_2",
				Value:  float64(rand.Int(1000)) / float64(rand.Int(20)),
				Source: "go_librato_example",
			},
			librato.Counter{
				Name:  "counter_1",
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
