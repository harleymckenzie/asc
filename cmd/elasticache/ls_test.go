package elasticache

import "testing"

func TestElasticacheColumns_ShowEndpoint(t *testing.T) {
	// Test when showEndpoint is true
	showEndpoint = true
	columns := elasticacheColumns()
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
	columns = elasticacheColumns()
	for _, col := range columns {
		if col.ID == "Endpoint" && col.Visible {
			t.Errorf("Expected Endpoint column to be hidden when showEndpoint is false")
		}
	}
}

func TestElasticacheColumns_SortName(t *testing.T) {
	// Test when sortName is true
	sortName = true
	columns := elasticacheColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Cache Name" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Cache Name column to be sortable when sortName is true")
	}

	// Test when sortName is false
	sortName = false
	columns = elasticacheColumns()
	for _, col := range columns {
		if col.ID == "Cache Name" && col.Sort {
			t.Errorf("Expected Cache Name column to not be sortable when sortName is false")
		}
	}
}
