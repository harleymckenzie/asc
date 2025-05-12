package smoke

import (
	"os"
	"os/exec"
	"testing"
)

// TestCloudFormationLsSmoke runs 'asc cloudformation ls' and prints the output for manual inspection.
func TestCloudFormationLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "cloudformation", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc cloudformation ls output:\n%s", out)
}

// TestCloudFormationLsSortByNameSmoke runs 'cloudformation ls -n' and prints the output for manual inspection.
func TestCloudFormationLsSortByNameSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "cloudformation", "ls", "-n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc cloudformation ls -n (sort by name) output:\n%s", out)
}

// TestCloudFormationLsSortByStatusSmoke runs 'cloudformation ls -s' and prints the output for manual inspection.
func TestCloudFormationLsSortByStatusSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "cloudformation", "ls", "-s")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc cloudformation ls -s (sort by status) output:\n%s", out)
}

// TestCloudFormationLsSortByLastUpdateSmoke runs 'cloudformation ls -u' and prints the output for manual inspection.
func TestCloudFormationLsSortByLastUpdateSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "cloudformation", "ls", "-u")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc cloudformation ls -u (sort by last update) output:\n%s", out)
}
