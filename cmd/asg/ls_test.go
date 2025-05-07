package asg

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/harleymckenzie/asc/pkg/service/asg"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

func TestAsgColumns_SortInstances(t *testing.T) {
	sortInstances = true
	columns := asgColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Instances" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Instances column to be sortable when sortInstances is true")
	}

	sortInstances = false
	columns = asgColumns()
	for _, col := range columns {
		if col.ID == "Instances" && col.Sort {
			t.Errorf("Expected Instances column to not be sortable when sortInstances is false")
		}
	}
}

func TestAsgColumns_SortDesired(t *testing.T) {
	sortDesiredCapacity = true
	columns := asgColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Desired" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Desired column to be sortable when sortDesiredCapacity is true")
	}

	sortDesiredCapacity = false
	columns = asgColumns()
	for _, col := range columns {
		if col.ID == "Desired" && col.Sort {
			t.Errorf("Expected Desired column to not be sortable when sortDesiredCapacity is false")
		}
	}
}

func TestAsgColumns_SortMin(t *testing.T) {
	sortMinCapacity = true
	columns := asgColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Min" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Min column to be sortable when sortMinCapacity is true")
	}

	sortMinCapacity = false
	columns = asgColumns()
	for _, col := range columns {
		if col.ID == "Min" && col.Sort {
			t.Errorf("Expected Min column to not be sortable when sortMinCapacity is false")
		}
	}
}

func TestAsgColumns_SortMax(t *testing.T) {
	sortMaxCapacity = true
	columns := asgColumns()
	found := false
	for _, col := range columns {
		if col.ID == "Max" && col.Sort {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected Max column to be sortable when sortMaxCapacity is true")
	}

	sortMaxCapacity = false
	columns = asgColumns()
	for _, col := range columns {
		if col.ID == "Max" && col.Sort {
			t.Errorf("Expected Max column to not be sortable when sortMaxCapacity is false")
		}
	}
}

func TestAsgSortingByMax(t *testing.T) {
	sortMaxCapacity = true
	groups := []types.AutoScalingGroup{
		{
			AutoScalingGroupName: strPtr("asg1"),
			MaxSize:              int32Ptr(1), // single digit
			MinSize:              int32Ptr(1), // single digit
			DesiredCapacity:      int32Ptr(1), // single digit
			Instances:            []types.Instance{},
		},
		{
			AutoScalingGroupName: strPtr("asg2"),
			MaxSize:              int32Ptr(16), // double digit
			MinSize:              int32Ptr(3),  // single digit
			DesiredCapacity:      int32Ptr(8),  // single digit
			Instances:            []types.Instance{},
		},
		{
			AutoScalingGroupName: strPtr("asg3"),
			MaxSize:              int32Ptr(2), // double digit
			MinSize:              int32Ptr(2), // double digit
			DesiredCapacity:      int32Ptr(13), // double digit
			Instances:            []types.Instance{},
		},
	}
	columns := asgColumns()
	selectedColumns, _ := tableformat.BuildColumns(columns)
	tbl := &asg.AutoScalingTable{
		AutoScalingGroups: groups,
		SelectedColumns:   selectedColumns,
	}

	headers := tbl.Headers()
	rows := tbl.Rows()

	// Find the index of the Max column
	maxIdx := -1
	for i, col := range selectedColumns {
		if col == "Max" {
			maxIdx = i
		}
	}
	if maxIdx == -1 {
		t.Fatalf("Max column not found in selectedColumns")
	}

	// Use go-pretty table.Writer to sort as in the CLI
	w := table.NewWriter()
	w.AppendHeader(headers)
	w.AppendRows(rows)
	w.SortBy([]table.SortBy{{Name: "Max", Mode: table.AscNumericAlpha}})

	var buf bytes.Buffer
	w.SetOutputMirror(&buf)
	w.RenderCSV()

	lines := strings.Split(buf.String(), "\n")
	// The first line is the header, skip it
	var maxVals []int
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		cols := strings.Split(line, ",")
		val := atoi(cols[maxIdx])
		maxVals = append(maxVals, val)
	}

	for i := 1; i < len(maxVals); i++ {
		if maxVals[i-1] > maxVals[i] {
			t.Errorf("Rows not sorted by Max ascending: %v", maxVals)
		}
	}
}

func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32 { return &i }
func atoi(s string) int {
	var n int
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}
