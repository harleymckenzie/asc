package smoke

import (
	"os"
	"os/exec"
	"testing"
)

// TestTargetGroupLsSmoke runs 'asc elb ls target-groups' and prints the output for manual inspection.
func TestTargetGroupLsSmoke(t *testing.T) {
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

