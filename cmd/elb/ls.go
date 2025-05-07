// The ls command list Elastic Load Balancers, as well as an alias for the relevant subcommand.
// It re-uses existing functions and flags from the relevant commands.

package elb

import (
	"context"
	"log"

	tg "github.com/harleymckenzie/asc/cmd/elb/target_group"
	"github.com/harleymckenzie/asc/pkg/service/elb"
	ascTypes "github.com/harleymckenzie/asc/pkg/service/elb/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/spf13/cobra"
)

var (
	list bool

	showARNs          bool
	showDNSName       bool
	showScheme        bool
	showAZs           bool
	showIPAddressType bool

	sortDNSName     bool
	sortType        bool
	sortCreatedTime bool
	sortScheme      bool
	sortVPCID       bool
)

func elbColumns() []tableformat.Column {
	return []tableformat.Column{
		{ID: "Name", Visible: true},
		{ID: "DNS Name", Visible: showDNSName, Sort: sortDNSName},
		{ID: "Scheme", Visible: showScheme, Sort: sortScheme},
		{ID: "State", Visible: true},
		{ID: "Type", Visible: true, Sort: sortType},
		{ID: "IP Type", Visible: showIPAddressType},
		{ID: "VPC ID", Visible: true, Sort: sortVPCID},
		{ID: "Created Time", Visible: true, Sort: sortCreatedTime},
		{ID: "ARN", Visible: showARNs, Sort: false},
		{ID: "Availability Zones", Visible: showAZs},
	}
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Elastic Load Balancers and target groups",
	Long: "List Elastic Load Balancers and target groups\n" +
		"  ls                           List all Elastic Load Balancers\n" +
		"  ls [elb-name]                List target groups for the specified ELB\n" +
		"  ls target-groups             List all target groups",
	GroupID: "actions",
	Run: func(cobraCmd *cobra.Command, args []string) {
		ctx := context.TODO()
		profile, _ := cobraCmd.Root().PersistentFlags().GetString("profile")
		region, _ := cobraCmd.Root().PersistentFlags().GetString("region")

		svc, err := elb.NewELBService(ctx, profile, region)
		if err != nil {
			log.Fatalf("Failed to initialize Elastic Load Balancing service: %v", err)
		}

		ListELBs(svc)
	},
}

// lsTargetGroupCmd calls the ListELBTargetGroups function, which is also used by the `target-group ls` command.
var lsTargetGroupCmd = &cobra.Command{
	Use:   "target-groups",
	Short: "List target groups",
	Run: func(cobraCmd *cobra.Command, args []string) {
		tg.ListTargetGroups(cobraCmd, args)
	},
}

func ListELBs(svc *elb.ELBService) {
	ctx := context.TODO()
	loadBalancers, err := svc.GetLoadBalancers(ctx, &ascTypes.GetLoadBalancersInput{})
	if err != nil {
		log.Fatalf("Failed to list Elastic Load Balancers: %v", err)
	}

	columns := elbColumns()
	selectedColumns, sortBy := tableformat.BuildColumns(columns)

	opts := tableformat.RenderOptions{
		SortBy: sortBy,
		List:   list,
		Title:  "Elastic Load Balancers",
	}

	tableformat.Render(&elb.ELBTable{
		LoadBalancers:   loadBalancers,
		SelectedColumns: selectedColumns,
	}, opts)
}

func addLsFlags(cobraCmd *cobra.Command) {
	// Output flags
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs Elastic Load Balancers in list format.")
	cobraCmd.Flags().BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each Elastic Load Balancer.")
	cobraCmd.Flags().BoolVarP(&showDNSName, "dns-name", "d", false, "Show the DNS name of the Elastic Load Balancer.")
	cobraCmd.Flags().BoolVarP(&showScheme, "scheme", "s", false, "Show the scheme for each Elastic Load Balancer.")
	cobraCmd.Flags().BoolVarP(&showAZs, "availability-zones", "z", false, "Show the availability zones for each Elastic Load Balancer.")
	cobraCmd.Flags().BoolVarP(&showIPAddressType, "ip-address-type", "i", false, "Show the IP address type for each Elastic Load Balancer.")

	// Sorting flags
	cobraCmd.Flags().BoolVarP(&sortDNSName, "sort-dns-name", "D", false, "Sort by descending DNS name.")
	cobraCmd.Flags().BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending Elastic Load Balancer type.")
	cobraCmd.Flags().BoolVarP(&sortCreatedTime, "sort-created-time", "t", false, "Sort by descending date created.")
	cobraCmd.Flags().BoolVarP(&sortScheme, "sort-scheme", "S", false, "Sort by descending scheme.")
	cobraCmd.Flags().BoolVarP(&sortVPCID, "sort-vpc-id", "V", false, "Sort by descending VPC ID.")
}

func init() {
	addLsFlags(lsCmd)

	lsCmd.AddCommand(lsTargetGroupCmd)
	tg.NewLsFlags(lsTargetGroupCmd)
}
