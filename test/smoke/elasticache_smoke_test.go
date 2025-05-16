package smoke

import (
	"os"
	"os/exec"
	"testing"
)

// TestElasticacheLsSmoke runs 'asc elasticache ls' and prints the output for manual inspection.
func TestElasticacheLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elasticache", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	t.Logf("asc elasticache ls output:\n%s", out)
}

// TestElasticacheLsSortByTypeSmoke runs 'elasticache ls -T' and prints the output for manual inspection.
func TestElasticacheLsSortByTypeSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elasticache", "ls", "-T")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	t.Logf("asc elasticache ls -T (sort by type) output:\n%s", out)
}

// TestElasticacheLsSortByStatusSmoke runs 'elasticache ls -s' and prints the output for manual inspection.
func TestElasticacheLsSortByStatusSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elasticache", "ls", "-s")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	t.Logf("asc elasticache ls -s (sort by status) output:\n%s", out)
}

// TestElasticacheLsSortByEngineSmoke runs 'elasticache ls -E' and prints the output for manual inspection.
func TestElasticacheLsSortByEngineSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elasticache", "ls", "-E")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
	t.Logf("asc elasticache ls -E (sort by engine) output:\n%s", out)
}
