package autolinks_test

import (
	"log"
	"os"

	autolinks "github.com/fundipper/goldmark-autolinks"
	"github.com/yuin/goldmark"
)

var source = []byte(`Standup notes:
- Previous day:
  - Gave feedback on TICKET-123.
  - Outlined presentation on syntax-aware markdown transformations.
  - Finished my part of TICKET-456 and assigned to Emily.
- Today:
  - Add integration tests for TICKET-789.
  - Create slides for presentation.`)

func Example() {
	md := goldmark.New(
		goldmark.WithExtensions(
			autolinks.NewExtender(
				map[string]string{
					`TICKET-\d+`: "https://example.com/TICKET?query=$0",
				}),
		),
	)

	if err := md.Convert(source, os.Stdout); err != nil {
		log.Fatal(err)
	}

	// Output:
	// <p>Standup notes:</p>
	// <ul>
	// <li>Previous day:
	// <ul>
	// <li>Gave feedback on <a href="https://example.com/TICKET?query=TICKET-123">TICKET-123</a>.</li>
	// <li>Outlined presentation on syntax-aware markdown transformations.</li>
	// <li>Finished my part of <a href="https://example.com/TICKET?query=TICKET-456">TICKET-456</a> and assigned to Emily.</li>
	// </ul>
	// </li>
	// <li>Today:
	// <ul>
	// <li>Add integration tests for <a href="https://example.com/TICKET?query=TICKET-789">TICKET-789</a>.</li>
	// <li>Create slides for presentation.</li>
	// </ul>
	// </li>
	// </ul>
}
