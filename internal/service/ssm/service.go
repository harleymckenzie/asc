package ssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

// SSMClientAPI is the interface for the SSM client.
type SSMClientAPI interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
	GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
	GetParameterHistory(ctx context.Context, params *ssm.GetParameterHistoryInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterHistoryOutput, error)
	PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
	DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
	DeleteParameters(ctx context.Context, params *ssm.DeleteParametersInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParametersOutput, error)
	DescribeParameters(ctx context.Context, params *ssm.DescribeParametersInput, optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
	LabelParameterVersion(ctx context.Context, params *ssm.LabelParameterVersionInput, optFns ...func(*ssm.Options)) (*ssm.LabelParameterVersionOutput, error)
	UnlabelParameterVersion(ctx context.Context, params *ssm.UnlabelParameterVersionInput, optFns ...func(*ssm.Options)) (*ssm.UnlabelParameterVersionOutput, error)
}

// SSMService is a struct that holds the SSM client.
type SSMService struct {
	Client SSMClientAPI
}

// NewSSMService creates a new SSM service.
func NewSSMService(ctx context.Context, profile string, region string) (*SSMService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}
	client := ssm.NewFromConfig(cfg.Config)

	return &SSMService{Client: client}, nil
}

// GetParameter fetches a single SSM parameter.
func (svc *SSMService) GetParameter(ctx context.Context, input *ascTypes.GetParameterInput) (*types.Parameter, error) {
	output, err := svc.Client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(input.Name),
		WithDecryption: aws.Bool(input.Decrypt),
	})
	if err != nil {
		return nil, err
	}
	if output.Parameter == nil {
		return nil, fmt.Errorf("parameter not found: %s", input.Name)
	}
	return output.Parameter, nil
}

// GetParameters fetches multiple SSM parameters by name.
func (svc *SSMService) GetParameters(ctx context.Context, input *ascTypes.GetParametersInput) ([]types.Parameter, error) {
	output, err := svc.Client.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          input.Names,
		WithDecryption: aws.Bool(input.Decrypt),
	})
	if err != nil {
		return nil, err
	}
	return output.Parameters, nil
}

// GetParametersByPath fetches all parameters under a path with pagination.
func (svc *SSMService) GetParametersByPath(ctx context.Context, input *ascTypes.GetParametersByPathInput) ([]types.Parameter, error) {
	var parameters []types.Parameter
	var nextToken *string

	for {
		output, err := svc.Client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:           aws.String(input.Path),
			Recursive:      aws.Bool(input.Recursive),
			WithDecryption: aws.Bool(input.Decrypt),
			NextToken:      nextToken,
		})
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, output.Parameters...)

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return parameters, nil
}

// DescribeParameters returns parameter metadata without values.
func (svc *SSMService) DescribeParameters(ctx context.Context, path string) ([]types.ParameterMetadata, error) {
	var parameters []types.ParameterMetadata
	var nextToken *string

	for {
		input := &ssm.DescribeParametersInput{
			NextToken: nextToken,
		}

		// Add path filter if provided
		if path != "" && path != "/" {
			input.ParameterFilters = []types.ParameterStringFilter{
				{
					Key:    aws.String("Path"),
					Option: aws.String("Recursive"),
					Values: []string{path},
				},
			}
		}

		output, err := svc.Client.DescribeParameters(ctx, input)
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, output.Parameters...)

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return parameters, nil
}

// PutParameter creates or updates a parameter.
func (svc *SSMService) PutParameter(ctx context.Context, input *ascTypes.PutParameterInput) error {
	putInput := &ssm.PutParameterInput{
		Name:      aws.String(input.Name),
		Value:     aws.String(input.Value),
		Type:      types.ParameterType(input.Type),
		Overwrite: aws.Bool(input.Overwrite),
	}

	if input.Description != "" {
		putInput.Description = aws.String(input.Description)
	}

	if len(input.Tags) > 0 {
		var tags []types.Tag
		for k, v := range input.Tags {
			tags = append(tags, types.Tag{
				Key:   aws.String(k),
				Value: aws.String(v),
			})
		}
		putInput.Tags = tags
	}

	_, err := svc.Client.PutParameter(ctx, putInput)
	return err
}

// DeleteParameter deletes a single parameter.
func (svc *SSMService) DeleteParameter(ctx context.Context, input *ascTypes.DeleteParameterInput) error {
	_, err := svc.Client.DeleteParameter(ctx, &ssm.DeleteParameterInput{
		Name: aws.String(input.Name),
	})
	return err
}

// DeleteParameters deletes multiple parameters (max 10 per call).
func (svc *SSMService) DeleteParameters(ctx context.Context, input *ascTypes.DeleteParametersInput) ([]string, error) {
	var failed []string

	// AWS allows max 10 parameters per DeleteParameters call
	for i := 0; i < len(input.Names); i += 10 {
		end := i + 10
		if end > len(input.Names) {
			end = len(input.Names)
		}

		batch := input.Names[i:end]
		output, err := svc.Client.DeleteParameters(ctx, &ssm.DeleteParametersInput{
			Names: batch,
		})
		if err != nil {
			return nil, err
		}

		failed = append(failed, output.InvalidParameters...)
	}

	return failed, nil
}

// CopyParameter copies a parameter to a new location.
func (svc *SSMService) CopyParameter(ctx context.Context, input *ascTypes.CopyParameterInput) error {
	// Validate inputs
	if input.Dest == "" {
		return fmt.Errorf("destination name cannot be empty")
	}
	if input.Source == input.Dest {
		return fmt.Errorf("source and destination cannot be the same")
	}

	// Get source parameter with decryption
	source, err := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
		Name:    input.Source,
		Decrypt: true,
	})
	if err != nil {
		return fmt.Errorf("get source parameter: %w", err)
	}

	// Put to destination
	err = svc.PutParameter(ctx, &ascTypes.PutParameterInput{
		Name:      input.Dest,
		Value:     aws.ToString(source.Value),
		Type:      string(source.Type),
		Overwrite: input.Overwrite,
	})
	if err != nil {
		return fmt.Errorf("put destination parameter: %w", err)
	}

	return nil
}

// CopyParametersRecursive copies all parameters from source path to destination path.
func (svc *SSMService) CopyParametersRecursive(ctx context.Context, sourcePath, destPath string, overwrite bool) (int, error) {
	// Get all parameters under source path
	params, err := svc.GetParametersByPath(ctx, &ascTypes.GetParametersByPathInput{
		Path:      sourcePath,
		Recursive: true,
		Decrypt:   true,
	})
	if err != nil {
		return 0, fmt.Errorf("get parameters by path: %w", err)
	}

	if len(params) == 0 {
		return 0, nil
	}

	copied := 0
	for _, param := range params {
		// Transform path
		newName := transformPath(sourcePath, destPath, aws.ToString(param.Name))

		err = svc.PutParameter(ctx, &ascTypes.PutParameterInput{
			Name:      newName,
			Value:     aws.ToString(param.Value),
			Type:      string(param.Type),
			Overwrite: overwrite,
		})
		if err != nil {
			return copied, fmt.Errorf("copy %s to %s: %w", aws.ToString(param.Name), newName, err)
		}
		copied++
	}

	return copied, nil
}

// MoveParameter moves a parameter to a new location (copy then delete).
func (svc *SSMService) MoveParameter(ctx context.Context, input *ascTypes.MoveParameterInput) error {
	// Copy first
	err := svc.CopyParameter(ctx, &ascTypes.CopyParameterInput{
		Source:    input.Source,
		Dest:      input.Dest,
		Overwrite: true,
	})
	if err != nil {
		return fmt.Errorf("copy parameter: %w", err)
	}

	// Delete source
	err = svc.DeleteParameter(ctx, &ascTypes.DeleteParameterInput{
		Name: input.Source,
	})
	if err != nil {
		return fmt.Errorf("delete source parameter: %w", err)
	}

	return nil
}

// MoveParametersRecursive moves all parameters from source path to destination path.
func (svc *SSMService) MoveParametersRecursive(ctx context.Context, sourcePath, destPath string) (int, error) {
	// Get all parameters under source path
	params, err := svc.GetParametersByPath(ctx, &ascTypes.GetParametersByPathInput{
		Path:      sourcePath,
		Recursive: true,
		Decrypt:   true,
	})
	if err != nil {
		return 0, fmt.Errorf("get parameters by path: %w", err)
	}

	if len(params) == 0 {
		return 0, nil
	}

	// First, copy all parameters to new location
	moved := 0
	var sourceNames []string
	var copiedNames []string
	for _, param := range params {
		name := aws.ToString(param.Name)
		sourceNames = append(sourceNames, name)

		// Transform path and copy
		newName := transformPath(sourcePath, destPath, name)
		err = svc.PutParameter(ctx, &ascTypes.PutParameterInput{
			Name:      newName,
			Value:     aws.ToString(param.Value),
			Type:      string(param.Type),
			Overwrite: true,
		})
		if err != nil {
			// Provide clear error message about partial state
			if len(copiedNames) > 0 {
				return moved, fmt.Errorf("copy %s to %s failed after copying %d parameters (destinations may need cleanup): %w",
					name, newName, moved, err)
			}
			return moved, fmt.Errorf("copy %s to %s: %w", name, newName, err)
		}
		copiedNames = append(copiedNames, newName)
		moved++
	}

	// Only delete source parameters after ALL copies succeed
	_, err = svc.DeleteParameters(ctx, &ascTypes.DeleteParametersInput{
		Names: sourceNames,
	})
	if err != nil {
		return moved, fmt.Errorf("delete source parameters (copies at destination exist): %w", err)
	}

	return moved, nil
}

// GetParameterHistory returns the version history of a parameter.
func (svc *SSMService) GetParameterHistory(ctx context.Context, input *ascTypes.GetParameterHistoryInput) ([]types.ParameterHistory, error) {
	var history []types.ParameterHistory
	var nextToken *string

	for {
		apiInput := &ssm.GetParameterHistoryInput{
			Name:           aws.String(input.Name),
			WithDecryption: aws.Bool(input.Decrypt),
			NextToken:      nextToken,
		}

		if input.MaxResults > 0 {
			apiInput.MaxResults = aws.Int32(int32(input.MaxResults))
		}

		output, err := svc.Client.GetParameterHistory(ctx, apiInput)
		if err != nil {
			return nil, err
		}

		history = append(history, output.Parameters...)

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken

		// If we have a max and reached it, stop
		if input.MaxResults > 0 && len(history) >= input.MaxResults {
			break
		}
	}

	return history, nil
}

// LabelParameterVersion adds labels to a specific parameter version.
func (svc *SSMService) LabelParameterVersion(ctx context.Context, input *ascTypes.LabelParameterVersionInput) ([]string, error) {
	apiInput := &ssm.LabelParameterVersionInput{
		Name:   aws.String(input.Name),
		Labels: input.Labels,
	}

	if input.Version > 0 {
		apiInput.ParameterVersion = aws.Int64(input.Version)
	}

	output, err := svc.Client.LabelParameterVersion(ctx, apiInput)
	if err != nil {
		return nil, err
	}

	return output.InvalidLabels, nil
}

// UnlabelParameterVersion removes labels from a parameter.
func (svc *SSMService) UnlabelParameterVersion(ctx context.Context, input *ascTypes.UnlabelParameterVersionInput) ([]string, error) {
	output, err := svc.Client.UnlabelParameterVersion(ctx, &ssm.UnlabelParameterVersionInput{
		Name:   aws.String(input.Name),
		Labels: input.Labels,
	})
	if err != nil {
		return nil, err
	}

	return output.InvalidLabels, nil
}

// transformPath converts /source/path/key to /dest/path/key.
func transformPath(sourcePath, destPath, paramName string) string {
	// Normalize paths - ensure no trailing slash for comparison
	sourcePath = strings.TrimSuffix(sourcePath, "/")
	destPath = strings.TrimSuffix(destPath, "/")

	// Replace source prefix with dest prefix
	return destPath + strings.TrimPrefix(paramName, sourcePath)
}
