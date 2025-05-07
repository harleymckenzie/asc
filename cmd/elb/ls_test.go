package elb

import "testing"

func TestElbColumns_ShowARNs(t *testing.T) {
	// Test when showARNs is true
	showARNs = true
	columns := elbColumns()
	found := false
	for _, col := range columns {
		if col.ID == "ARN" && col.Visible {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected ARN column to be visible when showARNs is true")
	}

	// Test when showARNs is false
	showARNs = false
	columns = elbColumns()
	for _, col := range columns {
		if col.ID == "ARN" && col.Visible {
			t.Errorf("Expected ARN column to be hidden when showARNs is false")
		}
	}
}

func TestElbColumns_SortDNSName(t *testing.T) {
	// Test when sortDNSName is true
	sortDNSName = true
	columns := elbColumns()
	found := false
	for _, col := range columns {
		if col.ID == "DNS Name" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected DNS Name column to be sortable when sortDNSName is true")
	}

	// Test when sortDNSName is false
	sortDNSName = false
	columns = elbColumns()
	for _, col := range columns {
		if col.ID == "DNS Name" && col.Sort {
			t.Errorf("Expected DNS Name column to not be sortable when sortDNSName is false")
		}
	}
}
