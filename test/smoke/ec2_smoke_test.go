package smoke

import (
	"os"
	"os/exec"
	"testing"
)

// TestEC2LsSmoke runs 'asc ec2 ls -L -t' and prints the output for manual inspection.
func TestEC2LsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "ls", "-L", "-t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 ls -L -t output:\n%s", out)
}

// TestEC2ShowSmoke runs 'asc ec2 show <instance-id>' and prints the output for manual inspection.
func TestEC2ShowSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	instanceID := os.Getenv("SMOKE_INSTANCE_ID")
	if instanceID == "" {
		t.Skip("set SMOKE_INSTANCE_ID to an instance ID to run this test")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "show", instanceID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 show %s output:\n%s", instanceID, out)
}
