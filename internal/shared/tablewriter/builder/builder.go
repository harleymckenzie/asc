package builder

import (
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
)

// Layout is the layout of a section.
type Layout string

// Layout constants.
const (
	Grid Layout = "grid"
	Row  Layout = "row"
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
// The Layout can be "Grid" or "Row"
type SectionConfig struct {
	Layout Layout
}

// AddSection uses the provided section to create the required tablewriter.HeaderRow and tablewriter.AttributeRow/tablewriter.HorizontalRow
// If the section uses the 'Grid' layout type, fields will be divided amongst AttributeRow's (num of Fields / t.columns)
func AddSection(t tablewriter.AscWriter, s Section) {
	// A section is made up of a header row and fields
	t.AppendHeaderRow(tablewriter.HeaderRow{
		Title: s.Title,
	})

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

			t.AppendAttributeRow(tablewriter.AttributeRow{
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

func BuildSections(fields []Field) []Section {
	categories := make(map[string][]tablewriter.Field)
	for _, f := range fields {
		categories[f.Category] = append(categories[f.Category], tablewriter.Field{
			Name:  f.Name,
			Value: f.Value,
		})
	}

	var sections []Section
	for name, fields := range categories {
		sections = append(sections, Section{
			Title:  name,
			Fields: fields,
			SectionConfig: SectionConfig{
				Layout: Grid,
			},
		})
	}
	return sections
}

func BuildTagSection(tags []tableformat.Tag) Section {
	var fields []tablewriter.Field
	for _, tag := range tags {
		fields = append(fields, tablewriter.Field{
			Name:  tag.Name,
			Value: tag.Value,
		})
	}

	section := Section{
		Title:  "Tags",
		Fields: fields,
		SectionConfig: SectionConfig{
			Layout: Row,
		},
	}
	return section
}
