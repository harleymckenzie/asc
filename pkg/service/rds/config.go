package rds

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
)

var DefaultTableConfig = tableformat.TableConfig{
	Columns: []tableformat.ColumnDefinition{
		{
			ID:        "cluster_identifier",
			Title:     "Cluster Identifier",
			AutoMerge: true,
		},
		{
			ID:       "identifier",
			Title:    "Identifier",
			WidthMax: 40,
		},
		{
			ID:       "status",
			Title:    "Status",
			WidthMax: 15,
		},
		{
			ID:       "engine",
			Title:    "Engine",
			WidthMax: 20,
		},
		{
			ID:       "size",
			Title:    "Size",
			WidthMax: 20,
		},
		{
			ID:       "role",
			Title:    "Role",
			WidthMax: 10,
		},
		{
			ID:       "endpoint",
			Title:    "Endpoint",
			WidthMax: 50,
		},
	},
	DefaultSort:  []string{"Cluster Identifier", "Identifier"},
	SeparateRows: true,
	MergeColumn:  aws.String("Cluster Identifier"),
}
