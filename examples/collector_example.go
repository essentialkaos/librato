package main

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"

	"pkg.re/essentialkaos/ek.v5/rand"

	"pkg.re/essentialkaos/librato.v3"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	librato.Mail = "mail@domain.com"
	librato.Token = "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234"

	collector := librato.NewCollector(time.Minute, collectSomeMetrics)
	collector.ErrorHandler = errorHandler

	for {
		time.Sleep(time.Hour)
	}
}

func collectSomeMetrics() []librato.Measurement {
	fmt.Println("Metrics collected")

	return []librato.Measurement{
		librato.Gauge{
			Name:  "example:gauge_1",
			Value: rand.Int(1000),
		},
		librato.Gauge{
			Name:   "example:gauge_2",
			Value:  float64(rand.Int(1000)) / float64(rand.Int(20)),
			Source: "go_librato_example",
		},
		librato.Counter{
			Name:  "example:counter_1",
			Value: rand.Int(1000),
		},
	}
}

func errorHandler(errs []error) {
	fmt.Println("Errors:")

	for _, err := range errs {
		fmt.Printf("  %v\n", err)
	}
}
