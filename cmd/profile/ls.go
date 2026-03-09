package profile

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/service/profile"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
	"github.com/harleymckenzie/asc/internal/shared/utils"
	"github.com/spf13/cobra"
)

var (
	list        bool
	showAll     bool
	showSSO     bool
	showRole    bool
	reverseSort bool
)

func init() {
	newLsFlags(lsCmd)
}

func getListFields() []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "Name", Category: "Profile", Visible: true, DefaultSort: true},
		{Name: "Type", Category: "Profile", Visible: showAll},
		{Name: "Region", Category: "Profile", Visible: true},
		{Name: "Output", Category: "Profile", Visible: true},
		{Name: "SSO Session", Category: "SSO", Visible: showSSO},
		{Name: "SSO Start URL", Category: "SSO", Visible: showSSO},
		{Name: "SSO Account ID", Category: "SSO", Visible: showSSO},
		{Name: "SSO Role Name", Category: "SSO", Visible: showSSO},
		{Name: "SSO Registration Scopes", Category: "SSO", Visible: showSSO},
		{Name: "Source Profile", Category: "Role", Visible: showRole},
		{Name: "Role ARN", Category: "Role", Visible: showRole},
	}
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List AWS CLI profiles",
	Aliases: []string{"list"},
	GroupID: "actions",
	Example: "  asc profile ls          # List all profiles\n" +
		"  asc profile ls -a       # Include SSO session entries\n" +
		"  asc profile ls --sso    # Show SSO details\n" +
		"  asc profile ls --role   # Show role assumption details",
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := profile.ListProfiles(profile.ListProfilesOptions{
			IncludeSSOSessions: showAll,
		})
		if err != nil {
			return fmt.Errorf("list profiles: %w", err)
		}

		tablewriter.RenderList(tablewriter.RenderListOptions{
			Title:         "Profiles",
			PlainStyle:    list,
			Fields:        getListFields(),
			Data:          utils.SlicesToAny(profiles),
			GetFieldValue: profile.GetFieldValue,
			GetTagValue:   profile.GetTagValue,
			ReverseSort:   reverseSort,
			HideEmpty:     true,
		})
		return nil
	},
}

func newLsFlags(cobraCmd *cobra.Command) {
	cobraCmd.Flags().BoolVarP(&list, "list", "l", false, "Output profiles in list format.")
	cobraCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Include SSO session entries.")
	cobraCmd.Flags().BoolVarP(&showSSO, "sso", "s", false, "Show SSO configuration details.")
	cobraCmd.Flags().BoolVarP(&showRole, "role", "R", false, "Show role assumption details.")
	cobraCmd.Flags().BoolVarP(&reverseSort, "reverse-sort", "r", false, "Reverse the sort order.")
}
