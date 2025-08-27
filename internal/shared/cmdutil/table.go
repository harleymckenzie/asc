package cmdutil

import (
	"fmt"

	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
)

// AttributeGetter is a function type that retrieves a field value from an instance
type AttributeGetter func(fieldName string, instance any) (string, error)

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
					fieldValue, err := getFieldValue(field.Name, instance)
					if err != nil {
						fieldValue = ""
					}
					instanceRow.Values = append(instanceRow.Values, fieldValue)
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

func PopulateFieldValues(instance any, fields []tablewriter.Field, getFieldValue AttributeGetter) ([]tablewriter.Field, error) {
	var populated []tablewriter.Field
	for _, field := range fields {
		if field.Category != "Tags" {
			fieldValue, err := getFieldValue(field.Name, instance)
			if err != nil {
				return nil, fmt.Errorf("get field value: %w", err)
			}
			populated = append(populated, tablewriter.Field{
				Category: field.Category,
				Name:     field.Name,
				Value:    fieldValue,
				Visible:  field.Visible,
			})
		} else {
			// For Tags category, keep the original field
			populated = append(populated, field)
		}
	}
	return populated, nil
}
