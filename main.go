package main

import (
	"github.com/harleymckenzie/asc/cmd"
	_ "github.com/harleymckenzie/asc/cmd/ec2"
	_ "github.com/harleymckenzie/asc/cmd/elasticache"
	_ "github.com/harleymckenzie/asc/cmd/rds"
)

func main() {
	cmd.Execute()
}
