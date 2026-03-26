package wait

import (
	"context"
	"fmt"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

// WaitHandler returns a status function and terminal check for a given resource.
// The status function polls the resource's current state. The terminal function
// returns true when the state is stable (no longer processing).
type WaitHandler func(ctx context.Context, profile, region string, uri *awsutil.ResourceURI) (
	statusFunc func(ctx context.Context) (string, error),
	isTerminal func(status string) bool,
	err error,
)

var handlers = map[string]WaitHandler{}

// RegisterHandler registers a wait handler for a service/resourceType combination.
// key format: "service/resourceType" (e.g. "ec2/instance", "rds/cluster").
func RegisterHandler(key string, handler WaitHandler) {
	handlers[key] = handler
}

// getHandler returns the registered handler for a ResourceURI, or an error if none exists.
func getHandler(uri *awsutil.ResourceURI) (WaitHandler, error) {
	key := fmt.Sprintf("%s/%s", uri.Service, uri.ResourceType)
	handler, ok := handlers[key]
	if !ok {
		return nil, fmt.Errorf("wait is not supported for %s %s", uri.Service, uri.ResourceType)
	}
	return handler, nil
}
