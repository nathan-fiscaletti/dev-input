# dev-input

[![Sponsor Me!](https://img.shields.io/badge/%F0%9F%92%B8-Sponsor%20Me!-blue)](https://github.com/sponsors/nathan-fiscaletti)
[![GoDoc](https://godoc.org/github.com/nathan-fiscaletti/dev-input?status.svg)](https://godoc.org/github.com/nathan-fiscaletti/dev-input)

This is a simple package for communicating with the Linux input filesystem.

## Features

- List input devices
- Filter input devices by event
  - Build in filters for keyboards, pointers and touch-pads
- Check if a device supports a particular event
- Check if a device supports a particular key-code
- Read events from input devices

## Usage

```sh
go get github.com/nathan-fiscaletti/dev-input
```

```go
package main

import "github.com/nathan-fiscaletti/dev-input"

func main() {
	keyboards, err := input.ListKeyboards()
	if err != nil {
		panic(err)
	}

	for _, keyboard := range keyboards {
		fmt.Printf("Keyboard: %s\n", keyboard.Name)
    }
}
```

See the [Go Docs](https://godoc.org/github.com/nathan-fiscaletti/dev-input) for more detailed usage.