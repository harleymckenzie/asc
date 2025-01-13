package ec2

import (
	"context"
	"log"

	"github.com/harleymckenzie/asc-go/cmd"
	"github.com/harleymckenzie/asc-go/pkg/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// EC2Cmd represents the ec2 command
var EC2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Perform EC2 operations",
}

// ec2 subcommands
var lsCmd = &cobra.Command{
	Use:              "ls",
	Short:            "List all EC2 instances",
	Long:             "List all EC2 instances. Sort flags can be combined (e.g. -iTn) to define multiple sort orders, where the order of the flags determines the sort priority.",
	TraverseChildren: true,
	PreRun: func(cobraCmd *cobra.Command, args []string) {
		// Clear any existing sort order
		sortOrder = []string{}

		// Visit flags in the order they appear in the command line
		cobraCmd.Flags().Visit(func(f *pflag.Flag) {
			switch f.Name {
			case "sort-name":
				sortOrder = append(sortOrder, "Name")
			case "sort-id":
				sortOrder = append(sortOrder, "Instance ID")
			case "sort-type":
				sortOrder = append(sortOrder, "Type")
			case "sort-created":
				sortOrder = append(sortOrder, "Created Time")
			}
		})
	},
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()

		svc, err := ec2.NewEC2Service(ctx, cmd.Profile)
		if err != nil {
			log.Fatalf("Failed to initialize EC2 service: %v", err)
		}

		err = svc.ListInstances(ctx, sortOrder, list)
		if err != nil {
			log.Fatalf("Error describing running instances: %v", err)
		}
	},
}

var (
	sortOrder []string
	list      bool
)

func init() {
	cmd.RootCmd.AddCommand(EC2Cmd)

	// ls sub command
	EC2Cmd.AddCommand(lsCmd)

	// Add flags - Options are displayed in the order they are added
	lsCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs EC2 instances in list format.")
	lsCmd.Flags().BoolP("sort-name", "n", true, "Sort by descending EC2 instance name.")
	lsCmd.Flags().BoolP("sort-id", "i", false, "Sort by descending EC2 instance Id.")
	lsCmd.Flags().BoolP("sort-type", "T", false, "Sort by descending EC2 instance type.")
	lsCmd.Flags().BoolP("sort-created", "t", false, "Sort by descending time created (most recently created first).")
	lsCmd.Flags().SortFlags = false
}
