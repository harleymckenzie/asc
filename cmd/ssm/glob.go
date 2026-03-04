package ssm

import (
	"context"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	ssmService "github.com/harleymckenzie/asc/internal/service/ssm"
)

// containsGlob returns true if the pattern contains glob characters.
func containsGlob(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[")
}

// resolveGlob resolves a glob pattern to matching SSM parameter names.
// If the pattern contains no glob characters, it returns the pattern as-is.
func resolveGlob(ctx context.Context, svc *ssmService.SSMService, pattern string) ([]string, error) {
	if !containsGlob(pattern) {
		return []string{pattern}, nil
	}

	// Extract the path prefix before the first glob character
	// e.g., "/myapp/prod/*" -> "/myapp/prod/"
	// e.g., "/myapp/*/key" -> "/myapp/"
	prefix := extractPrefix(pattern)

	// Fetch all parameters under the prefix
	metadata, err := svc.DescribeParameters(ctx, prefix)
	if err != nil {
		return nil, err
	}

	// Filter using path.Match
	var matches []string
	for _, m := range metadata {
		name := aws.ToString(m.Name)
		matched, err := path.Match(pattern, name)
		if err != nil {
			return nil, err
		}
		if matched {
			matches = append(matches, name)
		}
	}

	return matches, nil
}

// extractPrefix returns the path prefix before the first glob character.
// e.g., "/myapp/prod/*" -> "/myapp/prod/"
// e.g., "/myapp/*/key" -> "/myapp/"
// e.g., "/*/key" -> "/"
func extractPrefix(pattern string) string {
	// Find the first glob character
	idx := strings.IndexAny(pattern, "*?[")
	if idx == -1 {
		return pattern
	}

	// Find the last slash before the glob character
	prefix := pattern[:idx]
	if slashIdx := strings.LastIndex(prefix, "/"); slashIdx != -1 {
		return prefix[:slashIdx+1]
	}

	return "/"
}
