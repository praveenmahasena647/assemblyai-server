package main

import (
	"fmt"
	"os"

	"github.com/praveenmahasena647/server/api"
)

func main() {
	var err error = api.StartServer()

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
}
