package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
)

var (
	force bool
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop an EC2 instance",
	Aliases: []string{"shutdown", "halt"},
	Example: "asc ec2 stop i-1234567890abcdef0\n" +
		"asc ec2 stop i-1234567890abcdef0 --force",
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

		err = svc.StopInstance(ctx, &ascTypes.StopInstanceInput{
			InstanceID: args[0],
			Force:      force,
		})
		if err != nil {
			log.Fatalf("Failed to stop EC2 instance: %v", err)
		}

		ListEC2Instances(cobraCmd, []string{args[0]})
	},
}

func addStopFlags(stopCmd *cobra.Command) {
	stopCmd.Flags().BoolVarP(&force, "force", "f", false, "Force stop the EC2 instance")
}

func init() {
	addStopFlags(stopCmd)
}
