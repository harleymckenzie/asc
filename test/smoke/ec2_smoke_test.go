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
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
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
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
}

// TestEC2VolumeLsSmoke runs 'asc ec2 volume ls'
func TestEC2VolumeLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "volume", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 volume ls output:\n%s", out)
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
}

// TestEC2VolumeShowSmoke runs 'asc ec2 volume show <volume-id>'
func TestEC2VolumeShowSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	volumeID := os.Getenv("SMOKE_VOLUME_ID")
	if volumeID == "" {
		t.Skip("set SMOKE_VOLUME_ID to a volume ID to run this test")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "volume", "show", volumeID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 volume show %s output:\n%s", volumeID, out)
}

// TestEC2SnapshotLsSmoke runs 'asc ec2 snapshot ls'
func TestEC2SnapshotLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "snapshot", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 snapshot ls output:\n%s", out)
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
}

// TestEC2SnapshotShowSmoke runs 'asc ec2 snapshot show <snapshot-id>'
func TestEC2SnapshotShowSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	snapshotID := os.Getenv("SMOKE_SNAPSHOT_ID")
	if snapshotID == "" {
		t.Skip("set SMOKE_SNAPSHOT_ID to a snapshot ID to run this test")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "snapshot", "show", snapshotID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 snapshot show %s output:\n%s", snapshotID, out)
	if string(out) != "" && containsSortWarning(string(out)) {
		t.Fatalf("Test failed due to multiple sort fields: %s", out)
	}
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
}

// TestEC2AmiLsSmoke runs 'asc ec2 ami ls'
func TestEC2AmiLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "ami", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 ami ls output:\n%s", out)
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
}

// TestEC2AmiShowSmoke runs 'asc ec2 ami show <ami-id>'
func TestEC2AmiShowSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	amiID := os.Getenv("SMOKE_AMI_ID")
	if amiID == "" {
		t.Skip("set SMOKE_AMI_ID to an AMI ID to run this test")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "ami", "show", amiID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 ami show %s output:\n%s", amiID, out)
}

// TestEC2SecurityGroupLsSmoke runs 'asc ec2 security-group ls'
func TestEC2SecurityGroupLsSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "security-group", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 security-group ls output:\n%s", out)
	if containsAttributeError(string(out)) {
		t.Fatalf("Test failed due to attribute error: %s", out)
	}
	if containsMissingAttributeError(string(out)) {
		t.Fatalf("Test failed due to missing attribute error: %s", out)
	}
}

// TestEC2SecurityGroupShowSmoke runs 'asc ec2 security-group show <sg-id>'
func TestEC2SecurityGroupShowSmoke(t *testing.T) {
	if os.Getenv("SMOKE") != "1" {
		t.Skip("skipping smoke test; set SMOKE=1 to run")
	}
	sgID := os.Getenv("SMOKE_SG_ID")
	if sgID == "" {
		t.Skip("set SMOKE_SG_ID to a security group ID to run this test")
	}
	cmd := exec.Command("go", "run", "../../main.go", "ec2", "security-group", "show", sgID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\nOutput:\n%s", err, out)
	}
	t.Logf("asc ec2 security-group show %s output:\n%s", sgID, out)
}
