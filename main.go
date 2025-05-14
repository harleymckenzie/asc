package main

import (
	"fmt"
	"log"
	"os"

	"github.com/harleymckenzie/asc/cmd"
)

func main() {
	log.SetFlags(0)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
