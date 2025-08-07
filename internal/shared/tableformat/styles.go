package tableformat

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// TableStyle is a type alias for the table.Style type.
type TableStyle = table.Style

const DefaultTableStyle = "rounded"

// TableStyles is a map of table styles. It is used to select the style of the table.
var TableStyles = map[string]TableStyle{
	"list":              ListTableStyle,
	"rounded":           RoundedTableStyle,
	"rounded-separated": RoundedSeparatedTableStyle,
}

var (
	// ListTableStyle is a table style that displays a list of items.
	// NAME                INSTANCE ID          STATE    INSTANCE TYPE  PUBLIC IP
	// web-server-prod     i-0abc123def456789a  running  m7i.xlarge     18.134.56.201
	// app-server-staging  i-0def456789abc123b  running  r7i.xlarge     35.178.45.123
	// database-primary    i-0789abc123def456c  running  c7i.large      52.49.178.90
	ListTableStyle = table.Style{
		Name:    "list",
		Box:     table.StyleBoxRounded,
		Color:   table.ColorOptionsDefault,
		Format:  table.FormatOptionsDefault,
		Options: table.OptionsNoBordersAndSeparators,
	}

	// RoundedTableStyle is a table style that displays a list of items with rounded corners.
	// ╭──────────────────────────────────────────────────────────────────────────────────────╮
	// │ EC2 Instances                                                                        │
	// ├─────────────────────┬─────────────────────┬─────────┬───────────────┬────────────────┤
	// │ Name                │ Instance ID         │ State   │ Instance Type │ Public IP      │
	// ├─────────────────────┼─────────────────────┼─────────┼───────────────┼────────────────┤
	// │ web-server-prod     │ i-0abc123def456789a │ running │ m7i.xlarge    │ 18.134.56.201  │
	// │ app-server-staging  │ i-0def456789abc123b │ running │ r7i.xlarge    │ 35.178.45.123  │
	// │ database-primary    │ i-0789abc123def456c │ running │ c7i.large     │ 52.49.178.90   │
	// ╰─────────────────────┴─────────────────────┴─────────┴───────────────┴────────────────╯
	RoundedTableStyle = table.Style{
		Name:    "rounded",
		Box:     table.StyleBoxRounded,
		Color:   table.ColorOptions{Header: text.Colors{text.Bold}},
		Format:  table.FormatOptions{Header: text.FormatTitle},
		Options: table.OptionsDefault,
		Size:    table.SizeOptionsDefault,
		Title:   table.TitleOptions{Colors: text.Colors{text.Bold}},
	}

	// RoundedSeparatedTableStyle is a table style that displays a list of items with rounded corners and separated rows.
	// ╭──────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
	// │ RDS Clusters and Instances                                                                                   │
	// ├────────────────────┬────────────────────────────────────┬───────────┬──────────────┬────────────────┬────────┤
	// │ Cluster Identifier │ Instance Identifier                │ Status    │ Engine       │ Size           │ Role   │
	// ├────────────────────┼────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
	// │ prod-cluster-01    │ prod-cluster-01-instance-1         │ available │ aurora-mysql │ db.r6g.2xlarge │ Writer │
	// │                    ├────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
	// │                    │ prod-cluster-01-instance-2         │ available │ aurora-mysql │ db.r6g.2xlarge │ Reader │
	// ├────────────────────┼────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
	// │ dev-cluster-01     │ dev-cluster-01-instance-1          │ available │ aurora-mysql │ db.t3.medium   │ Reader │
	// │                    ├────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
	// │                    │ dev-cluster-01-instance-2          │ available │ aurora-mysql │ db.t4g.medium  │ Writer │
	// ╰────────────────────┴────────────────────────────────────┴───────────┴──────────────┴────────────────┴────────╯
	RoundedSeparatedTableStyle = table.Style{
		Name:    "rounded-separated",
		Box:     table.StyleBoxRounded,
		Color:   table.ColorOptions{Header: text.Colors{text.Bold}},
		Format:  table.FormatOptions{Header: text.FormatTitle},
		Options: DrawAllBordersAndSeparators,
		Size:    table.SizeOptionsDefault,
		Title:   table.TitleOptions{Colors: text.Colors{text.Bold}},
	}
)

// Options
var (
	DrawAllBordersAndSeparators = table.Options{
		DrawBorder:      true,
		SeparateColumns: true,
		SeparateFooter:  true,
		SeparateHeader:  true,
		SeparateRows:    true,
	}
)
