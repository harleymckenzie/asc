package tablewriter

import "fmt"

// ListTable handles simple list-style tables
type ListTable struct {
	Headers []string
	Rows    []Row
	Options AscTableRenderOptions
}

// TagGetter is a function type that retrieves a tag value from an instance
type TagGetter func(tagKey string, instance any) (string, error)

// NewListTable creates a new ListTable.
func NewListTable(options AscTableRenderOptions) *ListTable {
	return &ListTable{
		Options: options,
	}
}

func AppendTagFields(fields []Field, tags []string, instances []any) []Field {
	for _, tag := range tags {
		fields = append(fields, Field{Name: tag, Category: "Tags", Visible: true})
	}
	return fields
}

func (lt *ListTable) AddHeader(headers []string) {
	lt.Headers = headers
}

func (lt *ListTable) AddRow(row Row) {
	lt.Rows = append(lt.Rows, row)
}

func BuildHeaderRow(fields []Field) []string {
	headerRow := Row{
		Values: make([]string, 0, len(fields)),
	}
	for _, field := range fields {
		if field.Category == "Tags" {
			name := fmt.Sprintf("Tag: %s", field.Name)
			headerRow.Values = append(headerRow.Values, name)
			continue
		}
		if field.Visible {
			headerRow.Values = append(headerRow.Values, field.Name)
		}
	}
	return headerRow.Values
}

// BuildRows builds a slice of tablewriter.Row objects from a slice of instances and fields
// This is used to build the rows for a "List" table
func BuildRows(instances []any, fields []Field, getFieldValue AttributeGetter, getTagValue TagGetter) []Row {
	var rows []Row

	for _, instance := range instances {
		instanceRow := Row{Values: make([]string, 0, len(fields))}
		for _, field := range fields {
			if field.Category == "Tags" {
				fieldValue, err := getTagValue(field.Name, instance)
				if err != nil {
					fmt.Println("error getting tag value:", err)
					fieldValue = ""
				}
				instanceRow.Values = append(instanceRow.Values, fieldValue)
				continue
			}
			if field.Visible {
				fieldValue, err := getFieldValue(field.Name, instance)
				if err != nil {
					fmt.Println("error getting field value:", err)
					fieldValue = ""
				}
				instanceRow.Values = append(instanceRow.Values, fieldValue)
			}
		}
		rows = append(rows, instanceRow)
	}
	return rows
}

func (lt *ListTable) Render() {
	table := NewAscWriter(lt.Options)

	// Add headers
	if len(lt.Headers) > 0 {
		table.AppendHeader(lt.Headers)
	}

	// Add rows
	for _, row := range lt.Rows {
		table.AppendRow(row)
	}

	// Render the table
	table.Render()
}
