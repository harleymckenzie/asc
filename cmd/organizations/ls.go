package organizations

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/harleymckenzie/asc/internal/service/organizations"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/format"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	tree       bool
	showOUPath bool
)

func init() {
	newLsFlags(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List accounts and organizational units",
	Aliases: []string{"list"},
	GroupID: "actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(listOrganizations(cmd, args))
	},
}

func newLsFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&tree, "tree", "t", false, "Display as a tree of OUs and accounts")
	cmd.Flags().BoolVarP(&showOUPath, "ou-path", "P", false, "Show full OU path instead of direct parent OU")
}

func listOrganizations(cmd *cobra.Command, args []string) error {
	svc, err := cmdutil.CreateService(cmd, organizations.NewOrganizationsService)
	if err != nil {
		return fmt.Errorf("create new Organizations service: %w", err)
	}

	if tree {
		return listTree(cmd, svc)
	}
	return listTable(cmd, svc)
}

func listTable(cmd *cobra.Command, svc *organizations.OrganizationsService) error {
	accounts, err := svc.GetAccountsWithOUs(cmd.Context())
	if err != nil {
		return fmt.Errorf("list accounts: %w", err)
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Accounts",
		Style:         "rounded-separated",
		Fields:        organizations.AccountListFields(showOUPath),
		Data:          utils.SlicesToAny(accounts),
		GetFieldValue: organizations.GetFieldValue,
	})
	return nil
}

func listTree(cmd *cobra.Command, svc *organizations.OrganizationsService) error {
	roots, err := svc.GetRoots(cmd.Context())
	if err != nil {
		return fmt.Errorf("list roots: %w", err)
	}

	for _, root := range roots {
		tree, err := svc.BuildOUTree(cmd.Context(), aws.ToString(root.Id), aws.ToString(root.Name))
		if err != nil {
			return fmt.Errorf("build OU tree: %w", err)
		}
		printTree(tree, "", true)
	}
	return nil
}

// printTree recursively prints the OU tree with tree-drawing characters.
func printTree(node *organizations.OUTreeNode, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	if prefix == "" {
		fmt.Printf("%s (%s)\n", node.Name, node.ID)
	} else {
		fmt.Printf("%s%s%s (%s)\n", prefix, connector, node.Name, node.ID)
	}

	childPrefix := prefix
	if prefix != "" {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	totalChildren := len(node.Children) + len(node.Accounts)
	childIndex := 0

	for _, child := range node.Children {
		childIndex++
		childCopy := child
		printTree(&childCopy, childPrefix, childIndex == totalChildren)
	}

	for _, account := range node.Accounts {
		childIndex++
		isLastChild := childIndex == totalChildren
		accountConnector := "├── "
		if isLastChild {
			accountConnector = "└── "
		}
		accountStatus := formatAccountStatus(account)
		fmt.Printf("%s%s%s - %s%s\n", childPrefix, accountConnector, aws.ToString(account.Id), aws.ToString(account.Name), accountStatus)
	}
}

// formatAccountStatus returns a status suffix for suspended accounts.
func formatAccountStatus(account orgtypes.Account) string {
	if account.Status == orgtypes.AccountStatusSuspended {
		return " " + format.Status("suspended")
	}
	return ""
}
