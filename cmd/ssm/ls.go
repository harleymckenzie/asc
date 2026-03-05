package ssm

import (
	"context"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

// Variables
var (
	list        bool
	showValues  bool
	sortByName  bool
	sortByDate  bool
	reverseSort bool
)

// Init function
func init() {
	newLsFlags(lsCmd)
}

// getListFields returns a list of Field objects for displaying SSM parameters.
func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Parameter Details", Visible: true, DefaultSort: true, SortBy: sortByName, SortDirection: tablewriter.Asc},
		{Name: "Type", Category: "Parameter Details", Visible: true},
		{Name: "Value", Category: "Parameter Details", Visible: showValues},
		{Name: "Last Modified Date", Category: "Parameter Details", Visible: true, SortBy: sortByDate, SortDirection: tablewriter.Desc},
		{Name: "Last Modified User", Category: "Parameter Details", Visible: false},
		{Name: "Version", Category: "Parameter Details", Visible: true},
		{Name: "Tier", Category: "Parameter Details", Visible: false},
		{Name: "Description", Category: "Parameter Details", Visible: false},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls [path]",
	Short:   "List SSM parameters",
	Aliases: []string{"list"},
	GroupID: "actions",
	Example: "  asc ssm ls                    # List all parameters\n" +
		"  asc ssm ls /myapp/            # List parameters under path\n" +
		"  asc ssm ls /myapp/prod/       # List parameters under specific path",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(ListSSMParameters(cmd, args))
	},
}

// newLsFlags configures the flags for the ls command.
func newLsFlags(cobraCmd *cobra.Command) {
	// Output format flags
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Output parameters in list format.")
	cobraCmd.Flags().BoolVarP(&showValues, "values", "v", false, "Show parameter values in the output.")

	// Sorting flags
	cobraCmd.Flags().BoolVarP(&sortByName, "sort-name", "n", false, "Sort by parameter name.")
	cobraCmd.Flags().BoolVarP(&sortByDate, "sort-date", "d", false, "Sort by last modified date (most recent first).")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.MarkFlagsMutuallyExclusive("sort-name", "sort-date")
}

// ListSSMParameters lists SSM parameters with optional path filtering.
func ListSSMParameters(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	// Get path/pattern from args if provided
	pattern := ""
	if len(args) > 0 {
		pattern = args[0]
	}

	var resources []any
	isGlob := containsGlob(pattern)

	if showValues {
		if isGlob {
			// Resolve glob pattern to names, then batch-fetch values
			names, err := resolveGlob(ctx, svc, pattern)
			if err != nil {
				return fmt.Errorf("resolve glob: %w", err)
			}

			const batchSize = 10
			var allParameters []any
			for i := 0; i < len(names); i += batchSize {
				end := i + batchSize
				if end > len(names) {
					end = len(names)
				}
				parameters, err := svc.GetParameters(ctx, &ascTypes.GetParametersInput{
					Names:   names[i:end],
					Decrypt: true,
				})
				if err != nil {
					return fmt.Errorf("get parameters: %w", err)
				}
				allParameters = append(allParameters, utils.SlicesToAny(parameters)...)
			}
			resources = allParameters
		} else if pattern != "" {
			// Use GetParametersByPath when a non-glob path is provided
			parameters, err := svc.GetParametersByPath(ctx, &ascTypes.GetParametersByPathInput{
				Path:      pattern,
				Recursive: true,
				Decrypt:   true,
			})
			if err != nil {
				return fmt.Errorf("get parameters by path: %w", err)
			}
			resources = utils.SlicesToAny(parameters)
		} else {
			// No pattern: describe all, then batch-fetch values
			metadata, err := svc.DescribeParameters(ctx, pattern)
			if err != nil {
				return fmt.Errorf("describe parameters: %w", err)
			}

			names := make([]string, 0, len(metadata))
			for _, m := range metadata {
				names = append(names, aws.ToString(m.Name))
			}

			const batchSize = 10
			var allParameters []any
			for i := 0; i < len(names); i += batchSize {
				end := i + batchSize
				if end > len(names) {
					end = len(names)
				}
				parameters, err := svc.GetParameters(ctx, &ascTypes.GetParametersInput{
					Names:   names[i:end],
					Decrypt: true,
				})
				if err != nil {
					return fmt.Errorf("get parameters: %w", err)
				}
				allParameters = append(allParameters, utils.SlicesToAny(parameters)...)
			}
			resources = allParameters
		}
	} else {
		if isGlob {
			// Describe parameters under prefix, then filter by glob
			prefix := extractPrefix(pattern)
			metadata, err := svc.DescribeParameters(ctx, prefix)
			if err != nil {
				return fmt.Errorf("describe parameters: %w", err)
			}

			for _, m := range metadata {
				name := aws.ToString(m.Name)
				matched, err := path.Match(pattern, name)
				if err != nil {
					return fmt.Errorf("match pattern: %w", err)
				}
				if matched {
					resources = append(resources, m)
				}
			}
		} else {
			// Use DescribeParameters for listing (doesn't require decryption)
			parameters, err := svc.DescribeParameters(ctx, pattern)
			if err != nil {
				return fmt.Errorf("describe parameters: %w", err)
			}
			resources = utils.SlicesToAny(parameters)
		}
	}

	if len(resources) == 0 {
		if pattern != "" {
			fmt.Printf("No parameters found matching: %s\n", pattern)
		} else {
			fmt.Println("No parameters found.")
		}
		return nil
	}

	table := tablewriter.NewAscWriter(tablewriter.AscTableRenderOptions{
		Title: "Parameters",
	})
	if list {
		table.SetRenderStyle("plain")
	}

	fields := getListFields()
	headerRow := tablewriter.BuildHeaderRow(fields)
	table.AppendHeader(headerRow)
	table.AppendRows(tablewriter.BuildRows(resources, fields, ssm.GetFieldValue, ssm.GetTagValue))
	table.SetFieldConfigs(fields, reverseSort)

	table.Render()
	return nil
}
