package cloudformation

import "testing"

func TestCloudformationColumns_SortStatus(t *testing.T) {
	// Test when sortStatus is true
	sortStatus = true
	columns := cloudformationColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Status" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Status column to be sortable when sortStatus is true")
	}

	// Test when sortStatus is false
	sortStatus = false
	columns = cloudformationColumns()
	for _, col := range columns {
		if col.ID == "Status" && col.Sort {
			t.Errorf("Expected Status column to not be sortable when sortStatus is false")
		}
	}
}

func TestCloudformationColumns_SortName(t *testing.T) {
	// Test when sortName is true
	sortName = true
	columns := cloudformationColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Stack Name" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Stack Name column to be sortable when sortName is true")
	}

	// Test when sortName is false
	sortName = false
	columns = cloudformationColumns()
	for _, col := range columns {
		if col.ID == "Stack Name" && col.Sort {
			t.Errorf("Expected Stack Name column to not be sortable when sortName is false")
		}
	}
}
