package main

import (
	"github.com/paingha/stroll/app"
)

func main() {
	exampleApp := app.App{
		Name:        "example",
		Description: "an example application using the stroll cli",
		Version:     "v1.0",
	}
	if err := exampleApp.Run(); err != nil {
		panic(err)
	}
}
