// The ls command lists Elastic Load Balancers and target groups.

package elb

import (
	"context"
	"fmt"

	tg "github.com/harleymckenzie/asc/cmd/elb/target_group"
	"github.com/harleymckenzie/asc/internal/service/elb"
	ascTypes "github.com/harleymckenzie/asc/internal/service/elb/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
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

	reverseSort bool
)

// Init function
func init() {
	addLsFlags(lsCmd)

	lsCmd.AddCommand(lsTargetGroupCmd)
	tg.NewLsFlags(lsTargetGroupCmd)
}

// Column functions
func elbFields() []tableformat.Field {
	return []tableformat.Field{
		{ID: "Name", Display: true},
		{ID: "DNS Name", Display: showDNSName, Sort: sortDNSName},
		{ID: "Scheme", Display: showScheme, Sort: sortScheme},
		{ID: "State", Display: true},
		{ID: "Type", Display: true, Sort: sortType},
		{ID: "IP Type", Display: showIPAddressType},
		{ID: "VPC ID", Display: true, Sort: sortVPCID},
		{ID: "Created Time", Display: true, Sort: sortCreatedTime, SortDirection: "desc"},
		{ID: "ARN", Display: showARNs, Sort: false},
		{ID: "Availability Zones", Display: showAZs},
	}
}

// Command variable
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List Elastic Load Balancers and target groups",
	Long: "List Elastic Load Balancers and target groups\n" +
		"  ls                           List all Elastic Load Balancers\n" +
		"  ls [elb-name]                List target groups for the specified ELB\n" +
		"  ls target-groups             List all target groups",
	GroupID: "actions",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListELBs(cobraCmd, args))
	},
}

// Subcommand variable
var lsTargetGroupCmd = &cobra.Command{
	Use:   "target-groups",
	Short: "List target groups",
	Run: func(cobraCmd *cobra.Command, args []string) {
		tg.ListTargetGroups(cobraCmd, args)
	},
}

// Flag function
func addLsFlags(cobraCmd *cobra.Command) {
	// Output flags
	cobraCmd.Flags().
		BoolVarP(&list, "list", "l", false, "Outputs Elastic Load Balancers in list format.")
	cobraCmd.Flags().
		BoolVarP(&showARNs, "arn", "a", false, "Show ARNs for each Elastic Load Balancer.")
	cobraCmd.Flags().
		BoolVarP(&showDNSName, "dns-name", "d", false, "Show the DNS name of the Elastic Load Balancer.")
	cobraCmd.Flags().
		BoolVarP(&showScheme, "scheme", "s", false, "Show the scheme for each Elastic Load Balancer.")
	cobraCmd.Flags().
		BoolVarP(&showAZs, "availability-zones", "z", false, "Show the availability zones for each Elastic Load Balancer.")
	cobraCmd.Flags().
		BoolVarP(&showIPAddressType, "ip-address-type", "i", false, "Show the IP address type for each Elastic Load Balancer.")

	// Sorting flags
	cobraCmd.Flags().
		BoolVarP(&sortDNSName, "sort-dns-name", "D", false, "Sort by descending DNS name.")
	cobraCmd.Flags().
		BoolVarP(&sortType, "sort-type", "T", false, "Sort by descending Elastic Load Balancer type.")
	cobraCmd.Flags().
		BoolVarP(&sortCreatedTime, "sort-created-time", "t", false, "Sort by descending date created.")
	cobraCmd.Flags().BoolVarP(&sortScheme, "sort-scheme", "S", false, "Sort by descending scheme.")
	cobraCmd.Flags().BoolVarP(&sortVPCID, "sort-vpc-id", "V", false, "Sort by descending VPC ID.")
	cobraCmd.Flags().
		BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}

// Command functions
func ListELBs(cobraCmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	profile, region := cmdutil.GetPersistentFlags(cobraCmd)

	svc, err := elb.NewELBService(ctx, profile, region)
	if err != nil {
		return fmt.Errorf("create new Elastic Load Balancing service: %w", err)
	}

	loadBalancers, err := svc.GetLoadBalancers(ctx, &ascTypes.GetLoadBalancersInput{})
	if err != nil {
		return fmt.Errorf("list Elastic Load Balancers: %w", err)
	}

	fields := elbFields()

	opts := tableformat.RenderOptions{
		Title:  "Elastic Load Balancers",
		Style:  "rounded",
		SortBy: tableformat.GetSortByField(fields, reverseSort),
	}

	if list {
		opts.Style = "list"
	}

	tableformat.RenderTableList(&tableformat.ListTable{
		Instances: utils.SlicesToAny(loadBalancers),
		Fields:    fields,
		GetAttribute: func(fieldID string, instance any) (string, error) {
			return elb.GetAttributeValue(fieldID, instance)
		},
	}, opts)
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
