package main

import (
	"fmt"
	"glua"
	"os"
	"slices"
	"strings"
)

func main() {
	templates := []string{
		"simple",
		"content",
    "optional",
    "markdown",
	}
	if len(os.Args) == 1 || !slices.Contains(templates, os.Args[1]) {
		fmt.Println("Usage: go run ./examples/main.go {template}")
		fmt.Println("template can be one of:", strings.Join(templates, ", "))
		return
	}
	L := glua.New()
	err := L.Render(os.Stdout, nil, os.Args[1])
	if err != nil {
		panic(err)
	}
}
