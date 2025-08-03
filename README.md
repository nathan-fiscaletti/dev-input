# dev-input

[![Sponsor Me!](https://img.shields.io/badge/%F0%9F%92%B8-Sponsor%20Me!-blue)](https://github.com/sponsors/nathan-fiscaletti)
[![GoDoc](https://godoc.org/github.com/nathan-fiscaletti/dev-input?status.svg)](https://godoc.org/github.com/nathan-fiscaletti/dev-input)

This is a simple package for communicating with the Linux input filesystem.

## Features

- Read events from input devices
- List input devices
- Filter input devices by event
  - Built in filters for keyboards, pointers and touch-pads
- Check if a device supports a particular event
- Check if a device supports a particular key-code

## Usage

```sh
go get github.com/nathan-fiscaletti/dev-input
```

```go
package main

import (
	"context"
	"encoding/json"

	input "github.com/nathan-fiscaletti/dev-input"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	keyboards, err := input.ListKeyboards()
	if err != nil {
		panic(err)
	}

	for _, keyboard := range keyboards {
		go func(kb *input.Device) {
			err := kb.Listen(ctx, func(e input.Event) {
				printEvent(e)

				// Listen for escape
				if e.Code == 1 {
					cancel()
				}
			})
			if err != nil {
				panic(err)
			}
		}(keyboard)
	}

	<-ctx.Done()
}

func printEvent(v interface{}) {
	data, _ := json.Marshal(v)
	println(string(data))
}
```


See the [Go Docs](https://godoc.org/github.com/nathan-fiscaletti/dev-input) for more detailed usage.
