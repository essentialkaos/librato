package main

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"

	"pkg.re/essentialkaos/librato.v9"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	librato.Mail = "mail@domain.com"
	librato.Token = "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234"

	// Add annotation to example:annotation_1 stream
	errs := librato.AddAnnotation("example:annotation_1",
		librato.Annotation{
			Title:  "Deploy v31",
			Source: "server123",
			Desc:   "Revision: abcd1234",
			Links: []string{
				"https://build-service.com/build/31",
				"https://git-repo.com/commit/abcd1234",
			},
		},
	)

	if len(errs) != 0 {
		fmt.Println("Errors:")

		for _, err := range errs {
			fmt.Printf("  %v\n", err)
		}
	} else {
		fmt.Println("Annotation added")
	}

	time.Sleep(time.Minute)

	// Delete stream example:annotation_1 with all annotations
	errs = librato.DeleteAnnotations("example:annotation_1")

	if len(errs) != 0 {
		fmt.Println("Errors:")

		for _, err := range errs {
			fmt.Printf("  %v\n", err)
		}
	} else {
		fmt.Println("Annotation deleted")
	}
}
