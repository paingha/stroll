package main

import (
	"github.com/paingha/stroll"
)

func main() {
	exampleApp := stroll.App{
		Name:        "example",
		Description: "an example application using the stroll cli",
		Version:     "v1.0",
	}
	if err := exampleApp.Run(); err != nil {
		panic(err)
	}
}
