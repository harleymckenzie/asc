package tablewriter

// Section is a section of a table. It contains a title, fields, and a section config.
type Section struct {
	SectionConfig SectionConfig
	Fields        []Field
	Title         string
}

// SectionConfig is the configuration for a section.
// For now the SectionConfig only contains a "Layout"
// The Layout can be "Grid" or "Horizontal"
type SectionConfig struct {
	Layout Layout
}

// Layout is the layout of a section.
type Layout string

// Layout constants.
const (
	Grid       Layout = "Grid"
	Horizontal Layout = "Horizontal"
)
