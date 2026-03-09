package organizations

import "github.com/harleymckenzie/asc/internal/shared/tablewriter"

// AccountListFields returns the field definitions for the accounts list table.
func AccountListFields(showOUPath bool) []tablewriter.Field {
	return []tablewriter.Field{
		{Name: "OU", Category: "Organizations", Visible: !showOUPath, DefaultSort: true, Merge: true, SortDirection: tablewriter.Asc},
		{Name: "OU Path", Category: "Organizations", Visible: showOUPath, DefaultSort: true, Merge: true, SortDirection: tablewriter.Asc},
		{Name: "ID", Category: "Organizations", Visible: true},
		{Name: "Name", Category: "Organizations", Visible: true},
		{Name: "Email", Category: "Organizations", Visible: true},
		{Name: "Status", Category: "Organizations", Visible: true},
		{Name: "Joined Method", Category: "Organizations", Visible: true},
		{Name: "Joined", Category: "Organizations", Visible: true},
	}
}
