package smoke

import (
	"os"
	"os/exec"
	"testing"
)

// TestRDSLsSmoke runs 'asc rds ls -e' and prints the output for manual inspection.
func TestRDSLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls", "-e")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls -e output:\n%s", out)
}

// TestRDSLsBasicSmoke runs 'asc rds ls' and prints the output for manual inspection.
func TestRDSLsBasicSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls output:\n%s", out)
}

// TestRDSLsSortByNameSmoke runs 'rds ls -n' and prints the output for manual inspection.
func TestRDSLsSortByNameSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls", "-n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls -n (sort by name) output:\n%s", out)
}

// TestRDSLsSortByClusterSmoke runs 'rds ls -c' and prints the output for manual inspection.
func TestRDSLsSortByClusterSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls", "-c")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls -c (sort by cluster) output:\n%s", out)
}

// TestRDSLsSortByTypeSmoke runs 'rds ls -T' and prints the output for manual inspection.
func TestRDSLsSortByTypeSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls", "-T")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls -T (sort by type) output:\n%s", out)
}

// TestRDSLsSortByEngineSmoke runs 'rds ls -E' and prints the output for manual inspection.
func TestRDSLsSortByEngineSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls", "-E")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls -E (sort by engine) output:\n%s", out)
}

// TestRDSLsSortByStatusSmoke runs 'rds ls -s' and prints the output for manual inspection.
func TestRDSLsSortByStatusSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls", "-s")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls -s (sort by status) output:\n%s", out)
}

// TestRDSLsSortByRoleSmoke runs 'rds ls -R' and prints the output for manual inspection.
func TestRDSLsSortByRoleSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "rds", "ls", "-R")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	t.Logf("asc rds ls -R (sort by role) output:\n%s", out)
}
