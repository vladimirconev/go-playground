package main

import (
	"context"
	"os"

	"example.com/playground/cmd/app"
)

func main() {
	if err := app.New().RunContext(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
