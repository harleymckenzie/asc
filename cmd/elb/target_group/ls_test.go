package target_group

import "testing"

func TestTargetGroupColumns_ShowARNs(t *testing.T) {
	// Test when showARNs is true
	showARNs = true
	columns := targetGroupColumns()
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
	columns = targetGroupColumns()
	for _, col := range columns {
		if col.ID == "ARN" && col.Visible {
			t.Errorf("Expected ARN column to be hidden when showARNs is false")
		}
	}
}

func TestTargetGroupColumns_SortHealthCheckEnabled(t *testing.T) {
	// Test when showHealthCheckEnabled is true
	showHealthCheckEnabled = true
	columns := targetGroupColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Health Check Enabled" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Health Check Enabled column to be sortable when showHealthCheckEnabled is true")
	}

	// Test when showHealthCheckEnabled is false
	showHealthCheckEnabled = false
	columns = targetGroupColumns()
	for _, col := range columns {
		if col.ID == "Health Check Enabled" && col.Sort {
			t.Errorf("Expected Health Check Enabled column to not be sortable when showHealthCheckEnabled is false")
		}
	}
}
