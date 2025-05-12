package cmdutil

import (
	"errors"
	"log"
)

// DefaultErrorHandler handles AWS smithy.OperationError, context timeouts, and can be extended for other error types.
func DefaultErrorHandler(err error) error {
	if err == nil {
		return nil
	}
	// Print your context and the root cause
	var unwrapped error = err
	for {
		if next := errors.Unwrap(unwrapped); next != nil {
			unwrapped = next
		} else {
			break
		}
	}
	log.Fatalf("Error: %v\n", err) // This should include your prefix
	return nil              // Prevent Cobra from printing again
}
