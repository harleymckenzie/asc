package main

import (
	"github.com/harleymckenzie/asc/cmd"
	_ "github.com/harleymckenzie/asc/cmd/ec2"
	_ "github.com/harleymckenzie/asc/cmd/rds"
	_ "github.com/harleymckenzie/asc/cmd/elasticache"
)

var Version = "1.1.2"

func main() {
	cmd.NewRootCmd().Version = Version
	cmd.Execute()
}
