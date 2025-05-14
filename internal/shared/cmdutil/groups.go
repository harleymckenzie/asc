package cmdutil

import "github.com/spf13/cobra"


// ActionGroups returns the action groups.
func ActionGroups() []*cobra.Group {
	return []*cobra.Group{
		{ID: "actions", Title: "Action Commands"},
	}
}

// SubcommandGroups returns the subcommand groups.
func SubcommandGroups() []*cobra.Group {
	return []*cobra.Group{
		{ID: "subcommands", Title: "Subcommands"},
	}
}