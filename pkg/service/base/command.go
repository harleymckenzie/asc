package base

import (
	"context"
)

// CommandHandler defines the interface for handling service commands
type CommandHandler interface {
	Execute(ctx context.Context) error
}

// ResourceIdentifier represents a unique identifier for an AWS resource
type ResourceIdentifier struct {
	Name       string
	Type       string
	Additional map[string]string
}

// CommandOptions contains common options for all commands
type CommandOptions struct {
	Profile  string
	Region   string
	DryRun   bool
	Force    bool
	WaitSync bool
}

// ListOptions contains options specific to list commands
type ListOptions struct {
	CommandOptions
	SortOrder       []string
	List            bool
	SelectedColumns []string
}

// StateChangeOptions contains options for state change commands (start/stop/etc)
type StateChangeOptions struct {
	CommandOptions
	ResourceIDs []ResourceIdentifier
}

// ModifyOptions contains options for resource modification commands
type ModifyOptions struct {
	CommandOptions
	ResourceID    ResourceIdentifier
	Modifications map[string]interface{}
}
