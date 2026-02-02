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
	if len(opts.Tags) > 0 {
		fields = AppendTagFields(fields, opts.Tags, opts.Data)
	}

	table.AppendHeader(BuildHeaderRow(fields))
	table.AppendRows(BuildRows(opts.Data, fields, opts.GetFieldValue, opts.GetTagValue))
	table.SetFieldConfigs(fields, opts.ReverseSort)
	table.Render()
}
