package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc/pkg/service/ec2"
	"github.com/spf13/cobra"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/ec2/types"
)

var terminateCmd = &cobra.Command{
	Use:     "terminate",
	Short:   "Terminate an EC2 instance",
	Aliases: []string{"rm", "delete"},
	Example: "asc ec2 terminate i-1234567890abcdef0",
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

		err = svc.TerminateInstance(ctx, &ascTypes.TerminateInstanceInput{
			InstanceID: args[0],
		})
		if err != nil {
			log.Fatalf("Failed to terminate EC2 instance: %v", err)
		}

		ListEC2Instances(cobraCmd, ListInstancesInput{
			GetInstancesInput: &ascTypes.GetInstancesInput{
				InstanceIDs: []string{args[0]},
			},
			SelectedColumns: []string{"Name", "Instance ID", "State"},
		})
	},
}

func addTerminateFlags(terminateCmd *cobra.Command) {}

func init() {
	addTerminateFlags(terminateCmd)
}
