package cmdutil

import (
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
)

// AttributeGetter is a function type that retrieves a field value from an instance
type AttributeGetter func(fieldName string, instance any) string

// TagGetter is a function type that retrieves a tag value from an instance
type TagGetter func(tagKey string, instance any) (string, error)

func BuildHeaderRow(fields []tablewriter.Field) []string {
	headerRow := tablewriter.Row{
		Values: make([]string, 0, len(fields)),
	}
	for _, field := range fields {
		if field.Visible {
			headerRow.Values = append(headerRow.Values, field.Name)
		}
	}
	return headerRow.Values
}

// BuildRows builds a slice of tablewriter.Row objects from a slice of instances and fields
// This is used to build the rows for a "List" table
func BuildRows(instances []any, fields []tablewriter.Field, getFieldValue AttributeGetter, getTagValue TagGetter) []tablewriter.Row {
	var rows []tablewriter.Row

	for _, instance := range instances {
		instanceRow := tablewriter.Row{
			Values: make([]string, 0, len(fields)),
		}
		for _, field := range fields {
			if field.Visible {
				if field.Category == "Tags" {
					fieldValue, err := getTagValue(field.Name, instance)
					if err != nil {
						fieldValue = ""
					}
					instanceRow.Values = append(instanceRow.Values, fieldValue)
				} else {
					instanceRow.Values = append(instanceRow.Values, getFieldValue(field.Name, instance))
				}
			}
		}
		rows = append(rows, instanceRow)
	}
	return rows
}

func AppendTagFields(fields []tablewriter.Field, tags []string, instances []any) []tablewriter.Field {
	for _, tag := range tags {
		fields = append(fields, tablewriter.Field{Name: tag, Category: "Tags", Visible: true})
	}
	return fields
}
