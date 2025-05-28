package tableformat

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

// GetSortByField returns a table.SortBy for the given fields
// It will look at the 'Sort', 'DefaultSort' and 'SortDirection' fields to determine the sort order
// If no sort fields are found, the 'DefaultSort' field will be used
// If no 'DefaultSort' field is found, the first field will be used
func GetSortByField(fields []Field, reverseSort bool) []table.SortBy {
	if len(fields) == 0 {
		return nil
	}

	var sortBy []table.SortBy
	// 1. Find fields with Sort == true
	for _, f := range fields {
		if f.Sort {
			sortBy = append(sortBy, table.SortBy{
				Name:       f.ID,
				Mode:       sortDirection(f.SortDirection, reverseSort),
				IgnoreCase: true,
			})
		}
	}

	// 2. If no sort fields found, find fields with DefaultSort == true
	// If no sort direction is specified, use ascending order (and reverse if specified)
	if len(sortBy) == 0 {
		for _, f := range fields {
			if f.DefaultSort {
				if f.SortDirection == "" {
					sortBy = append(sortBy, table.SortBy{
						Name:       f.ID,
						Mode:       sortDirection("asc", reverseSort),
						IgnoreCase: true,
					})
				} else {
					sortBy = append(sortBy, table.SortBy{
						Name:       f.ID,
						Mode:       sortDirection(f.SortDirection, reverseSort),
						IgnoreCase: true,
					})
				}
			}
		}
	}

	// 3. If still none, use the first field
	if len(sortBy) == 0 {
		sortBy = append(sortBy, table.SortBy{
			Name:       fields[0].ID,
			Mode:       sortDirection("asc", reverseSort),
			IgnoreCase: true,
		})
	}

	if len(sortBy) > 1 {
		fmt.Printf(
			"Warning: Multiple sort fields found: %v\nFields will be sorted by the first field\n",
			sortBy,
		)
	}

	return sortBy
}

// sortDirection returns the sort mode for the given direction
// If reverseSort is true, the direction will be reversed
func sortDirection(direction string, reverseSort bool) table.SortMode {
	if reverseSort {
		switch direction {
		case "asc":
			return table.DscNumericAlpha
		case "desc", "":
			return table.AscNumericAlpha
		}
	}

	if direction == "desc" {
		return table.DscNumericAlpha
	}
	return table.AscNumericAlpha
}
