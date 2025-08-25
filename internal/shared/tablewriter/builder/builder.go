package builder

import (
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
)

// Layout is the layout of a section.
type Layout string

// Layout constants.
const (
	Grid       Layout = "grid"
	Horizontal Layout = "horizontal"
)

type Field struct {
	Category string
	Name     string
	Value    string
	Visible  bool
}

// Section is a section of a table. It contains a title, fields, and a section config.
type Section struct {
	SectionConfig SectionConfig
	Fields        []tablewriter.Field
	Title         string
}

// SectionConfig is the configuration for a section.
// For now the SectionConfig only contains a "Layout"
// The Layout can be "Grid" or "Horizontal"
type SectionConfig struct {
	Layout Layout
}

// AddSection uses the provided section to create the required tablewriter.HeaderRow and tablewriter.AttributeRow/tablewriter.HorizontalRow
// If the section uses the 'Grid' layout type, fields will be divided amongst AttributeRow's (num of Fields / t.columns)
func AddSection(t tablewriter.AscWriter, s Section) {
	// A section is made up of a header row and fields
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

			t.AppendGridRow(tablewriter.GridRow{
				Fields: s.Fields[start:end],
			})
		}
	} else {
		for _, f := range s.Fields {
			t.AppendHorizontalRow(tablewriter.HorizontalRow{
				Field: f,
			})
		}
	}
}

// AddSections accepts a list of Sections and creates new sections for each.
func AddSections(t tablewriter.AscWriter, s []Section) {
	for _, sec := range s {
		AddSection(t, sec)
	}
}

// BuildSections builds a list of sections from a list of fields with a grid layout.
// The title of each section is the category of the fields.
func BuildSections(fields []Field, layout Layout) []Section {
	categories := make(map[string][]tablewriter.Field)
	categoryOrder := make([]string, 0)

	for _, f := range fields {
		if _, exists := categories[f.Category]; !exists {
			categoryOrder = append(categoryOrder, f.Category)
		}
		categories[f.Category] = append(categories[f.Category], tablewriter.Field{
			Name:  f.Name,
			Value: f.Value,
		})
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

// BuildSection builds a section with a row layout.
func BuildSection(title string, fields []Field, layout Layout) Section {
	var rowFields []tablewriter.Field
	for _, f := range fields {
		rowFields = append(rowFields, tablewriter.Field{
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
