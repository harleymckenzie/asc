// The ls command lists Elastic Load Balancers and target groups.

package elb

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/elb"
	ascTypes "github.com/harleymckenzie/asc/internal/service/elb/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
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
}

// Column functions
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Load Balancer Details", Visible: true, DefaultSort: true},
		{Name: "DNS Name", Category: "Network", Visible: showDNSName, SortBy: sortDNSName, SortDirection: tablewriter.Asc},
		{Name: "Scheme", Category: "Load Balancer Details", Visible: showScheme, SortBy: sortScheme, SortDirection: tablewriter.Asc},
		{Name: "State", Category: "Load Balancer Details", Visible: true},
		{Name: "Type", Category: "Load Balancer Details", Visible: true, SortBy: sortType, SortDirection: tablewriter.Asc},
		{Name: "IP Type", Category: "Network", Visible: showIPAddressType},
		{Name: "VPC ID", Category: "Network", Visible: true, SortBy: sortVPCID, SortDirection: tablewriter.Asc},
		{Name: "Created Time", Category: "Load Balancer Details", Visible: true, SortBy: sortCreatedTime, SortDirection: tablewriter.Desc},
		{Name: "ARN", Category: "Load Balancer Details", Visible: showARNs},
		{Name: "Availability Zones", Category: "Network", Visible: showAZs},
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

// Flag function
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
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}

// Command functions
func ListELBs(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, elb.NewELBService)
	if err != nil {
		return fmt.Errorf("create elb service: %w", err)
	}

	loadBalancers, err := svc.GetLoadBalancers(cmd.Context(), &ascTypes.GetLoadBalancersInput{})
	if err != nil {
		return fmt.Errorf("get load balancers: %w", err)
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Elastic Load Balancers",
	})
	if list {
		table.SetStyle("plain")
	}
	fields := getListFields()

	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(loadBalancers), fields, elb.GetFieldValue, elb.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}
