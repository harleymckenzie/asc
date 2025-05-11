package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
)

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start an EC2 instance",
	Example: "asc ec2 start i-1234567890abcdef0",
	GroupID: "actions",
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

		// Start the instance
		err = svc.StartInstance(ctx, &ascTypes.StartInstanceInput{
			InstanceID: args[0],
		})
		if err != nil {
			log.Fatalf("Failed to start EC2 instance: %v", err)
		}

		ListEC2Instances(cobraCmd, []string{args[0]})
	},
}

func newStartFlags(cobraCmd *cobra.Command) {}

func init() {
	newStartFlags(startCmd)
}
