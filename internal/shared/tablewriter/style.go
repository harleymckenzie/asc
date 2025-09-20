package tablewriter

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var TableStyles = map[string]table.Style{
	"rounded":   StyleRounded,
	"plain":     StylePlain,
	"separated": StyleRoundedSeparated,
}

// SetStyle sets the style of the table.
// SetStyle sets the style of the table
func (at *AscTable) SetStyle(style string) {
	switch style {
	case "rounded":
		at.table.SetStyle(table.StyleRounded)
	case "plain":
		at.table.SetStyle(StylePlain)
	case "separated":
		at.table.SetStyle(StyleRoundedSeparated)
	default:
		at.table.SetStyle(table.StyleRounded)
	}
}

var (
	// StyleRounded is a table style that displays a list of items with rounded corners.
	// ╭──────────────────────────────────────────────────────────────────────────────────────╮
	// │ EC2 Instances                                                                        │
	// ├─────────────────────┬─────────────────────┬─────────┬───────────────┬────────────────┤
	// │ Name                │ Instance ID         │ State   │ Instance Type │ Public IP      │
	// ├─────────────────────┼─────────────────────┼─────────┼───────────────┼────────────────┤
	// │ web-server-prod     │ i-0abc123def456789a │ running │ m7i.xlarge    │ 18.134.56.201  │
	// │ app-server-staging  │ i-0def456789abc123b │ running │ r7i.xlarge    │ 35.178.45.123  │
	// │ database-primary    │ i-0789abc123def456c │ running │ c7i.large     │ 52.49.178.90   │
	// ╰─────────────────────┴─────────────────────┴─────────┴───────────────┴────────────────╯
	StyleRounded = table.Style{
		Name:    "rounded",
		Box:     table.StyleBoxRounded,
		Color:   ColorOptionsDefault,
		Format:  table.FormatOptions{Header: text.FormatTitle},
		Options: table.OptionsDefault,
		Size:    table.SizeOptionsDefault,
		Title:   table.TitleOptionsDefault,
	}

	// StylePlain is a table style that displays a list of items with no borders or separators.
	// NAME                INSTANCE ID          STATE    INSTANCE TYPE  PUBLIC IP
	// web-server-prod     i-0abc123def456789a  running  m7i.xlarge     18.134.56.201
	// app-server-staging  i-0def456789abc123b  running  r7i.xlarge     35.178.45.123
	// database-primary    i-0789abc123def456c  running  c7i.large      52.49.178.90
	StylePlain = table.Style{
		Name:    "list",
		Box:     table.StyleBoxRounded,
		Color:   ColorOptionsDefault,
		Format:  table.FormatOptionsDefault,
		Options: table.OptionsNoBordersAndSeparators,
	}

	// Style RoundedSeparated is a table style that displays a list of items with rounded corners and separated rows.
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
	StyleRoundedSeparated = table.Style{
		Name:    "rounded-separated",
		Box:     table.StyleBoxRounded,
		Color:   ColorOptionsDefault,
		Format:  table.FormatOptions{Header: text.FormatTitle},
		Options: DrawAllBordersAndSeparators,
		Size:    table.SizeOptionsDefault,
		Title:   table.TitleOptionsDefault,
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

// ColorOptions
var (
	ColorOptionsDefault = table.ColorOptions{
		Border:       text.Colors{},          // borders (if nil, uses one of the below)
		Footer:       text.Colors{},          // footer row(s) colors
		Header:       text.Colors{text.Bold}, // header row(s) colors
		IndexColumn:  text.Colors{},          // index-column colors (row #, etc.)
		Row:          text.Colors{},          // regular row(s) colors
		RowAlternate: text.Colors{},          // regular row(s) colors for the even-numbered rows
		Separator:    text.Colors{},          // separators (if nil, uses one of the above)
	}
)
