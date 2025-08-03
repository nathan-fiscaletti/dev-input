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
