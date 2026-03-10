package taskdefinition

import (
	"fmt"
	"strings"

	"github.com/harleymckenzie/asc/internal/service/ecs"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	reverseSort bool
)

func init() {
	newLsFlags(lsCmd)
}

func getFamilyListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Family", Category: "Task Definition", Visible: true, DefaultSort: true},
	}
}

func getRevisionListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Family", Category: "Task Definition", Visible: true},
		{Name: "Revision", Category: "Task Definition", Visible: true, DefaultSort: true},
		{Name: "ARN", Category: "Task Definition", Visible: true},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls [family-name]",
	Short:   "List task definition families, or revisions for a specific family",
	Long:    "Without arguments, lists all task definition families. With a family name argument, lists all revisions for that family.",
	Aliases: []string{"list"},
	GroupID: "actions",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return cmdutil.DefaultErrorHandler(ListTaskDefinitionRevisions(cmd, args[0]))
		}
		return cmdutil.DefaultErrorHandler(ListTaskDefinitionFamilies(cmd))
	},
}

func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Outputs in list format.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
	cobraCmd.Flags().SortFlags = false
}

func ListTaskDefinitionFamilies(cmd *cobra.Command) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	families, err := svc.ListTaskDefinitionFamilies(cmd.Context(), &ascTypes.ListTaskDefinitionFamiliesInput{})
	if err != nil {
		return fmt.Errorf("list task definition families: %w", err)
	}

	// Convert to TaskDefinitionFamily structs
	var data []ecs.TaskDefinitionFamily
	for _, f := range families {
		data = append(data, ecs.TaskDefinitionFamily{Name: f})
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         "Task Definition Families",
		PlainStyle:    list,
		Fields:        getFamilyListFields(),
		Data:          utils.SlicesToAny(data),
		GetFieldValue: ecs.GetFieldValue,
		GetTagValue:   ecs.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}

func ListTaskDefinitionRevisions(cmd *cobra.Command, familyName string) error {
	svc, err := cmdutil.CreateService(cmd, ecs.NewECSService)
	if err != nil {
		return fmt.Errorf("create ECS service: %w", err)
	}

	arns, err := svc.ListTaskDefinitionRevisions(cmd.Context(), &ascTypes.ListTaskDefinitionRevisionsInput{
		FamilyName: familyName,
	})
	if err != nil {
		return fmt.Errorf("list task definition revisions: %w", err)
	}

	// Convert ARNs to TaskDefinitionRevision structs
	var data []ecs.TaskDefinitionRevision
	for _, arn := range arns {
		shortName := ecs.ShortARN(arn)
		parts := strings.SplitN(shortName, ":", 2)
		family := shortName
		revision := ""
		if len(parts) == 2 {
			family = parts[0]
			revision = parts[1]
		}
		data = append(data, ecs.TaskDefinitionRevision{
			ARN:      arn,
			Family:   family,
			Revision: revision,
		})
	}

	tablewriter.RenderList(tablewriter.RenderListOptions{
		Title:         fmt.Sprintf("Task Definition Revisions (%s)", familyName),
		PlainStyle:    list,
		Fields:        getRevisionListFields(),
		Data:          utils.SlicesToAny(data),
		GetFieldValue: ecs.GetFieldValue,
		GetTagValue:   ecs.GetTagValue,
		ReverseSort:   reverseSort,
	})
	return nil
}
