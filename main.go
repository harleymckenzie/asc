package main

import (
	"github.com/harleymckenzie/asc-go/cmd"
	_ "github.com/harleymckenzie/asc-go/cmd/ec2"
	_ "github.com/harleymckenzie/asc-go/cmd/rds"
	_ "github.com/harleymckenzie/asc-go/cmd/elasticache"
)

var Version = "1.1.1"

func main() {
	cmd.RootCmd.Version = Version
	cmd.Execute()
}
