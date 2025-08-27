package tablewriter

// ListTable handles simple list-style tables
type ListTable struct {
	Headers []string
	Rows    []Row
	Options AscTableRenderOptions
}

// NewListTable creates a new ListTable.
func NewListTable(options AscTableRenderOptions) *ListTable {
	return &ListTable{
		Options: options,
	}
}

func (lt *ListTable) AddHeader(headers []string) {
	lt.Headers = headers
}

func (lt *ListTable) AddRow(row Row) {
	lt.Rows = append(lt.Rows, row)
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
