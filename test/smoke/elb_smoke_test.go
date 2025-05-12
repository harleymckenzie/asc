package smoke

import (
	"os"
	"os/exec"
	"testing"
)

// TestELBLsSmoke runs 'asc elb ls' and prints the output for manual inspection.
func TestELBLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elb", "ls")
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
	t.Logf("asc elb ls output:\n%s", out)
}

// TestELBLsTargetGroupsSmoke runs 'asc elb ls target-groups' and prints the output for manual inspection.
func TestELBLsTargetGroupsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elb", "ls", "target-groups")
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
	t.Logf("asc elb ls target-groups output:\n%s", out)
}

// TestELBLsSortByDNSNameSmoke runs 'elb ls -D' and prints the output for manual inspection.
func TestELBLsSortByDNSNameSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elb", "ls", "-D")
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
	t.Logf("asc elb ls -D (sort by DNS name) output:\n%s", out)
}

// TestELBLsSortByTypeSmoke runs 'elb ls -T' and prints the output for manual inspection.
func TestELBLsSortByTypeSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elb", "ls", "-T")
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
	t.Logf("asc elb ls -T (sort by type) output:\n%s", out)
}

// TestELBLsSortByCreatedTimeSmoke runs 'elb ls -t' and prints the output for manual inspection.
func TestELBLsSortByCreatedTimeSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elb", "ls", "-t")
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
	t.Logf("asc elb ls -t (sort by created time) output:\n%s", out)
}

// TestELBLsSortBySchemeSmoke runs 'elb ls -S' and prints the output for manual inspection.
func TestELBLsSortBySchemeSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elb", "ls", "-s", "-S")
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
	t.Logf("asc elb ls -S (sort by scheme) output:\n%s", out)
}

// TestELBLsSortByVPCIDSmoke runs 'elb ls -V' and prints the output for manual inspection.
func TestELBLsSortByVPCIDSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "elb", "ls", "-V")
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
	t.Logf("asc elb ls -V (sort by VPC ID) output:\n%s", out)
}
