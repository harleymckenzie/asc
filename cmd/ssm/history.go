package ssm

import (
	"fmt"
	"strings"

	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	historyDecrypt bool
	historyLimit   int
	historyList    bool
)

func init() {
	newHistoryFlags(historyCmd)
}

func getHistoryFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Version", Category: "Version Details", Visible: true, DefaultSort: true, SortDirection: tablewriter.Desc},
		{Name: "Value", Category: "Version Details", Visible: true},
		{Name: "Type", Category: "Version Details", Visible: true},
		{Name: "Last Modified Date", Category: "Version Details", Visible: true},
		{Name: "Last Modified User", Category: "Version Details", Visible: false},
		{Name: "Labels", Category: "Version Details", Visible: true},
		{Name: "Description", Category: "Version Details", Visible: false},
	}
}

var historyCmd = &cobra.Command{
	Use:     "history <parameter-name>",
	Short:   "Show version history of an SSM parameter",
	GroupID: "actions",
	Args:    cobra.ExactArgs(1),
	Example: `  asc ssm history /myapp/prod/config
  asc ssm history /myapp/prod/secret --decrypt
  asc ssm history /myapp/prod/config --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ShowParameterHistory(cmd, args))
	},
}

func newHistoryFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&historyDecrypt, "decrypt", "d", false, "Decrypt and show SecureString values")
	cmd.Flags().IntVarP(&historyLimit, "limit", "n", 0, "Limit number of versions to show (0 = all)")
	cmd.Flags().BoolVarP(&historyList, "list", "l", false, "Output in list format")
}

func ShowParameterHistory(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	paramName := args[0]

	history, err := svc.GetParameterHistory(ctx, &ascTypes.GetParameterHistoryInput{
		Name:       paramName,
		Decrypt:    historyDecrypt,
		MaxResults: historyLimit,
	})
	if err != nil {
		return fmt.Errorf("get parameter history: %w", err)
	}

	if len(history) == 0 {
		fmt.Printf("No history found for parameter: %s\n", paramName)
		return nil
	}

	// Extract just the parameter name for the title (without path)
	displayName := paramName
	if idx := strings.LastIndex(paramName, "/"); idx != -1 {
		displayName = paramName[idx+1:]
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: fmt.Sprintf("History: %s", displayName),
	})
	if historyList {
		table.SetRenderStyle("plain")
	}

	fields := getHistoryFields()
	table.AppendHeader(tablewriter.BuildHeaderRow(fields))
	table.AppendRows(tablewriter.BuildRows(utils.SlicesToAny(history), fields, ssm.GetFieldValue, ssm.GetTagValue))
	table.SetFieldConfigs(fields, false)

	table.Render()
	return nil
}
