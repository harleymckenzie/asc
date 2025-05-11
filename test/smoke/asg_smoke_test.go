package smoke

import (
	"os"
	"os/exec"
	"testing"
)

// TestASGLsSmoke runs 'asc asg ls' and prints the output for manual inspection.
func TestASGLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "asg", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc asg ls output:\n%s", out)
}

// TestASGLsSchedulesSmoke runs 'asc asg ls schedules' and prints the output for manual inspection.
func TestASGLsSchedulesSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "asg", "ls", "schedules")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc asg ls schedules output:\n%s", out)
}

// TestASGLsNameSmoke runs 'asc asg ls <asg-name>' if ASG_SMOKE_NAME is set.
func TestASGLsNameSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	asgName := os.Getenv("ASG_SMOKE_NAME")
	if asgName == "" {
		t.Skip("set ASG_SMOKE_NAME to an ASG name to run this test")
	}
	cmd := exec.Command("go", "run", "../../main.go", "asg", "ls", asgName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc asg ls %s output:\n%s", asgName, out)
}

// TestASGLsSortByNameSmoke runs 'asg ls -n' and prints the output for manual inspection.
func TestASGLsSortByNameSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "asg", "ls", "-n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc asg ls -n (sort by name) output:\n%s", out)
}

// TestASGLsSortByDesiredCapacitySmoke runs 'asg ls -d' and prints the output for manual inspection.
func TestASGLsSortByDesiredCapacitySmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "asg", "ls", "-d")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc asg ls -d (sort by desired capacity) output:\n%s", out)
}

// TestASGLsSortByMinCapacitySmoke runs 'asg ls -m' and prints the output for manual inspection.
func TestASGLsSortByMinCapacitySmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "asg", "ls", "-m")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc asg ls -m (sort by min capacity) output:\n%s", out)
}

// TestASGLsSortByMaxCapacitySmoke runs 'asg ls -M' and prints the output for manual inspection.
func TestASGLsSortByMaxCapacitySmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "asg", "ls", "-M")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc asg ls -M (sort by max capacity) output:\n%s", out)
}
