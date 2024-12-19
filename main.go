package main

import (
	"github.com/harleymckenzie/asc-go/cmd"
	_ "github.com/harleymckenzie/asc-go/cmd/ec2"
	_ "github.com/harleymckenzie/asc-go/cmd/rds"
)

func main() {
	cmd.Execute()
}
