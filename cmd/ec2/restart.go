package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
)

var restartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Restart an EC2 instance",
	Aliases: []string{"reboot"},
	Example: "asc ec2 restart i-1234567890abcdef0",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		// If an argument hasn't been provided, print the help message
		if len(args) == 0 {
			cobraCmd.Help()
			return
		}

		svc, err := ec2.NewEC2Service(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize EC2 service: %v", err)
		}

		err = svc.RestartInstance(ctx, &ascTypes.RestartInstanceInput{
			InstanceID: args[0],
		})
		if err != nil {
			log.Fatalf("Failed to restart EC2 instance: %v", err)
		}

		fmt.Printf("Reboot request sent to instance %s\n", args[0])
	},
}

func addRestartFlags(restartCmd *cobra.Command) {}

func init() {
	addRestartFlags(restartCmd)
}
