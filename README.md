# UnixID

A Go library for generating unique, time-based IDs using Unix timestamps at nanosecond precision.

## Overview

UnixID provides functionality for generating and managing unique identifiers with the following features:

- High-performance ID generation based on Unix nanosecond timestamps
- Thread-safe concurrent ID generation
- Built-in collision avoidance through sequential numbering
- Support for both server-side and client-side (WebAssembly) environments
- Date conversion utilities for timestamp-to-date formatting

## Installation

```bash
go get github.com/cdvelop/unixid
```

## Quick Start

### Server-side Usage

```go
package main

import (
	"fmt"
	"github.com/cdvelop/unixid"
)

func main() {
	// Create a new UnixID handler (server-side)
	idHandler, err := unixid.NewUnixID()
	if err != nil {
		panic(err)
	}

	// Generate a new unique ID
	id, err := idHandler.GetNewID()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated ID: %s\n", id)
	// Output: Generated ID: 1624397134562544800

	// Convert an ID to a human-readable date
	dateStr, err := idHandler.UnixNanoToStringDate(id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID timestamp represents: %s\n", dateStr)
	// Output: ID timestamp represents: 2021-06-23 15:38:54
}
```

### Client-side (WebAssembly) Usage

For WebAssembly environments, you need to provide a session number handler:

```go
// Example session handler implementation
type sessionHandler struct{}

func (sessionHandler) userSessionNumber() (number string, err error) {
	// In a real application, this would return the user's session number
	return "42", nil
}

// Create a new UnixID handler with session handler
idHandler, err := unixid.NewUnixID(&sessionHandler{})
```

## ID Format

The generated IDs follow this format:

- Server-side: `[unix_timestamp_in_nanoseconds]` (e.g., `1624397134562544800`)
- Client-side: `[unix_timestamp_in_nanoseconds].[user_session_number]` (e.g., `1624397134562544800.42`)

## API Reference

### Core Functions

- `NewUnixID(...)`: Creates a new UnixID handler for ID generation
- `GetNewID()`: Generates a new unique ID
- `UnixNanoToStringDate(unixNanoStr)`: Converts a Unix nanosecond timestamp ID to a human-readable date

### Additional Utility Functions

- `UnixSecondsToTime(unixSeconds any) string`: Converts a Unix timestamp in seconds to a formatted time string (HH:mm:ss). e.g., `1624397134` -> `15:38:54` supports `int64`, `string`, and `float64` types


## Validate ID

The `ValidateID` function validates and parses a given ID string. It returns the parsed ID as an `int64` and an error if the ID is invalid.

### Example

```go
package main

import (
	"fmt"
	"github.com/cdvelop/unixid"
)

func main() {
	id := "1624397134562544800"
	parsedID, err := unixid.ValidateID(id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed ID: %d\n", parsedID)
	// Output: Parsed ID: 1624397134562544800
}
```

## Thread Safety

The library handles concurrent ID generation safely through mutex locking in server-side environments.

## License

See the [LICENSE](LICENSE) file for details.