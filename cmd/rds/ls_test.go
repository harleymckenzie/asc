package rds

import "testing"

func TestRdsColumns_ShowEndpoint(t *testing.T) {
	// Test when showEndpoint is true
	showEndpoint = true
	columns := rdsColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Endpoint" && col.Visible {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Endpoint column to be visible when showEndpoint is true")
	}

	// Test when showEndpoint is false
	showEndpoint = false
	columns = rdsColumns()
	for _, col := range columns {
		if col.ID == "Endpoint" && col.Visible {
			t.Errorf("Expected Endpoint column to be hidden when showEndpoint is false")
		}
	}
}

func TestRdsColumns_SortName(t *testing.T) {
	// Test when sortName is true
	sortName = true
	columns := rdsColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Identifier" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Identifier column to be sortable when sortName is true")
	}

	// Test when sortName is false
	sortName = false
	columns = rdsColumns()
	for _, col := range columns {
		if col.ID == "Identifier" && col.Sort {
			t.Errorf("Expected Identifier column to not be sortable when sortName is false")
		}
	}
}
