### Librato [![GoDoc](https://godoc.org/pkg.re/essentialkaos/librato.v1?status.svg)](https://godoc.org/pkg.re/essentialkaos/librato.v1)

Package for working with [Librato Metrics](https://www.librato.com) API from go code.

#### Installation

````
go get pkg.re/essentialkaos/librato.v1
````

#### Status

This package is unnder heavy construction, please do not use in production code.

#### Example

```Go
package main

// ////////////////////////////////////////////////////////////////////////////////// //

import (
  "os"
  "time"

  "pkg.re/essentialkaos/librato.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {

  // Set auth credentials which will be used for all actions
  librato.Mail = "mail@domain.com"
  librato.Token = "mysupertokenhere"

  var errs []error

  // Add annotation to service:annotation stream
  errs = librato.AddAnnotation("service:annotation",
    &librato.Annotation{
      Title:  "Deploy v31",
      Source: "server123",
      Desc:   "Revision: abcd1234",
      Links: []string{
        "https://build-service.com/build/31",
        "https://git-repo.com/commit/abcd1234",
      },
    },
  )

  // Exit with 1 if we have errors
  if len(errs) != 0 {
    os.Exit(1)
  }

  // Delete stream service:annotation with all annotations
  errs = librato.DeleteAnnotations("service:annotation")

  if len(errs) != 0 {
    os.Exit(1)
  }

  // Add counter
  errs = librato.AddMetric(
    &librato.Counter{
      Name:  "service:random1",
      Value: 345,
    },
  )

  if len(errs) != 0 {
    os.Exit(1)
  }

  // Add gauge
  errs = librato.AddMetric(
    &librato.Gauge{
      Name:  "service:random2",
      Value: 45.2,
    },
  )

  if len(errs) != 0 {
    os.Exit(1)
  }

  // Create struct for async sending metrics data
  // With this preferences metrics will be sended to Librato once
  // a minute or if queue size reached 60 elements
  metrics, err := librato.NewMetrics(time.Minute, 60)

  if err != nil {
    os.Exit(1)
  }

  metrics.Add(
    &librato.Counter{
      Name:  "service:random1",
      Value: 345,
    },
  )

  metrics.Add(
    &librato.Gauge{
      Name:  "service:random2",
      Value: 45.2,
    },
  )

  // Force sending metrics before exit
  metrics.Send()
}
```

#### License

[EKOL](https://essentialkaos.com/ekol)
