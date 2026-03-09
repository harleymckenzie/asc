package tablewriter

// RenderListOptions contains the configuration for rendering a list table.
type RenderListOptions struct {
	Title         string
	Style         string
	PlainStyle    bool
	Fields        []Field
	Tags          []string
	Data          []any
	GetFieldValue AttributeGetter
	GetTagValue   TagGetter
	ReverseSort   bool
	HideEmpty     bool
}

// hideEmptyFields sets Visible to false for any visible field where all data values are empty.
func hideEmptyFields(fields []Field, data []any, getFieldValue AttributeGetter) []Field {
	for i := range fields {
		if !fields[i].Visible {
			continue
		}
		allEmpty := true
		for _, item := range data {
			val, err := getFieldValue(fields[i].Name, item)
			if err == nil && val != "" {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			fields[i].Visible = false
		}
	}
	return fields
}

// SetFieldVisibility sets the visibility of a field by name.
// It iterates through the fields slice and sets Visible to the specified value
// for the first field matching the given name.
func SetFieldVisibility(fields []Field, name string, visible bool) {
	for i := range fields {
		if fields[i].Name == name {
			fields[i].Visible = visible
			return
		}
	}
}

// RenderList creates and renders a list-style table with the provided options.
// This helper consolidates the common pattern used across all list commands:
//   - Creates a table with the specified title and style
//   - Applies plain style if PlainStyle is true
//   - Appends tag fields to the field list
//   - Builds and appends header row
//   - Builds and appends data rows using the provided getters
//   - Configures field sorting
//   - Renders the table
func RenderList(opts RenderListOptions) {
	table := NewAscWriter(AscTableRenderOptions{
		Title: opts.Title,
		Style: opts.Style,
	})

	if opts.PlainStyle {
		table.SetRenderStyle("plain")
	}

	fields := opts.Fields
	if opts.HideEmpty {
		fields = hideEmptyFields(fields, opts.Data, opts.GetFieldValue)
	}
	if len(opts.Tags) > 0 {
		fields = AppendTagFields(fields, opts.Tags, opts.Data)
	}

	table.AppendHeader(BuildHeaderRow(fields))
	table.AppendRows(BuildRows(opts.Data, fields, opts.GetFieldValue, opts.GetTagValue))
	table.SetFieldConfigs(fields, opts.ReverseSort)
	table.Render()
}
