package tablewriter

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
)

// DetailTable handles detailed tables.
type DetailTable struct {
	Headers  []string
	Rows     []Row
	Options  AscTableRenderOptions
	Sections []Section
}

// NewDetailTable creates a new DetailTable.
func NewDetailTable(options AscTableRenderOptions) *DetailTable {
	return &DetailTable{
		Options: options,
	}
}

// AddSection adds a section to the detail table
func (dt *DetailTable) AddSection(s Section) {
	dt.Sections = append(dt.Sections, s)
}

// AddSections adds multiple sections to the detail table
func (dt *DetailTable) AddSections(sections []Section) {
	dt.Sections = append(dt.Sections, sections...)
}

// AddSections adds multiple sections to a table writer
func AddSections(t AscWriter, sections []Section) {
	for _, section := range sections {
		AddSection(t, section)
	}
}

// AddSection adds a single section to a table writer
func AddSection(t AscWriter, s Section) {
	// Add section title
	t.AppendTitleRow(s.Title)

	if s.SectionConfig.Layout == Grid {
		numFields := len(s.Fields)
		columns := t.GetColumns()

		// Calculate number of rows needed, ensuring at least 1 row if there are fields
		numRows := (numFields + columns - 1) / columns // Ceiling division

		for i := 0; i < numRows; i++ {
			start := i * columns
			end := start + columns
			if end > numFields {
				end = numFields
			}

			t.AppendGridRow(GridRow{
				Fields: s.Fields[start:end],
			})
		}
	} else {
		for _, f := range s.Fields {
			t.AppendHorizontalRow(HorizontalRow{
				Field: f,
			})
		}
	}
}

// BuildSection builds a section with a row layout.
func BuildSection(title string, fields []Field, layout Layout) Section {
	var rowFields []Field
	for _, f := range fields {
		rowFields = append(rowFields, Field{
			Name:  f.Name,
			Value: f.Value,
		})
	}

	section := Section{
		Title:  title,
		Fields: rowFields,
		SectionConfig: SectionConfig{
			Layout: layout,
		},
	}
	return section
}

// BuildSections builds a list of sections from a list of fields with a grid layout.
// The title of each section is the category of the fields.
func BuildSections(fields []Field, layout Layout) []Section {
	categories := make(map[string][]Field)
	categoryOrder := make([]string, 0)

	for _, f := range fields {
		if f.Visible {
			if _, exists := categories[f.Category]; !exists {
				categoryOrder = append(categoryOrder, f.Category)
			}
			categories[f.Category] = append(categories[f.Category], Field{
				Name:  f.Name,
				Value: f.Value,
			})
		}
	}

	var sections []Section
	for _, name := range categoryOrder {
		sections = append(sections, Section{
			Title:  name,
			Fields: categories[name],
			SectionConfig: SectionConfig{
				Layout: layout,
			},
		})
	}
	return sections
}

func (dt *DetailTable) AddHeader(headers []string) {
	dt.Headers = headers
}

func (dt *DetailTable) AddRow(row Row) {
	dt.Rows = append(dt.Rows, row)
}

// SetColumnWidth sets the minimum and maximum width for all columns.
func (at *AscTable) SetColumnWidth(minWidth int, maxWidth int) {
	// TODO: Implement this for DetailTable only
	configs := make([]table.ColumnConfig, at.renderOptions.Columns)
	for i := 0; i < at.renderOptions.Columns; i++ {
		configs[i] = table.ColumnConfig{Number: i + 1, WidthMin: minWidth, WidthMax: maxWidth}
	}
	at.table.SetColumnConfigs(configs)
}

func (dt *DetailTable) Render() {
	table := NewAscWriter(dt.Options)

	// Add headers
	if len(dt.Headers) > 0 {
		table.AppendHeader(dt.Headers)
	}

	// Add rows
	for _, row := range dt.Rows {
		table.AppendRow(row)
	}

	// Add sections with layout logic
	for _, section := range dt.Sections {
		// Add section title
		table.AppendTitleRow(section.Title)

		if section.SectionConfig.Layout == Grid {
			numFields := len(section.Fields)
			columns := table.GetColumns()

			// Calculate number of rows needed, ensuring at least 1 row if there are fields
			numRows := (numFields + columns - 1) / columns // Ceiling division

			for i := 0; i < numRows; i++ {
				start := i * columns
				end := start + columns
				if end > numFields {
					end = numFields
				}

				table.AppendGridRow(GridRow{
					Fields: section.Fields[start:end],
				})
			}
		} else {
			for _, f := range section.Fields {
				table.AppendHorizontalRow(HorizontalRow{
					Field: f,
				})
			}
		}
	}

	// Render the table
	table.Render()
}

func PopulateFieldValues(instance any, fields []Field, getFieldValue AttributeGetter) ([]Field, error) {
	var populated []Field
	for _, field := range fields {
		if field.Category != "Tags" {
			fieldValue, err := getFieldValue(field.Name, instance)
			if err != nil {
				return nil, fmt.Errorf("get field value: %w", err)
			}
			populated = append(populated, Field{
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