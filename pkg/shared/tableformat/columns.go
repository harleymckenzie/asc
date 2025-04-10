package tableformat

type ColumnDefinition struct {
	ID        string
	Title     string
	WidthMin  int
	WidthMax  int
	AutoMerge bool
}

type TableConfig struct {
	Columns      []ColumnDefinition
	DefaultSort  []string
	SeparateRows bool
	MergeColumn  *string
}

// GetSelectedColumns returns the column headers for the selected column IDs
func GetSelectedColumns(config TableConfig, selectedIDs []string) []string {
	headers := make([]string, 0)
	for _, id := range selectedIDs {
		for _, col := range config.Columns {
			if col.ID == id {
				headers = append(headers, col.Title)
				break
			}
		}
	}
	return headers
}
