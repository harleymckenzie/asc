package cloudformation

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Integration test for NewCloudFormationService (skipped unless INTEGRATION=1)
func TestNewCloudFormationService_Integration(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("skipping integration test; set INTEGRATION=1 to run")
	}
	svc, err := NewCloudFormationService(context.Background(), "", "eu-west-1")
	assert.NoError(t, err)
	assert.NotNil(t, svc)
}
