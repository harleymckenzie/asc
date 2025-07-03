package iam

import (
	"context"

	iam "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	ascTypes "github.com/harleymckenzie/asc/internal/service/iam/types"
	"github.com/harleymckenzie/asc/internal/shared/awsutil"
	"github.com/harleymckenzie/asc/internal/shared/utils"
)

type IAMClientAPI interface {
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
	ListRolePolicies(ctx context.Context, params *iam.ListRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListRolePoliciesOutput, error)
	ListRoles(ctx context.Context, params *iam.ListRolesInput, optFns ...func(*iam.Options)) (*iam.ListRolesOutput, error)
	ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
}

type IAMService struct {
	Client IAMClientAPI
}

//
// Service functions
//

// NewIAMService creates a new IAM service.
func NewIAMService(ctx context.Context, profile string, region string) (*IAMService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}
	client := iam.NewFromConfig(cfg.Config)

	return &IAMService{Client: client}, nil
}

// GetRoles gets IAM roles. If a role name is provided, it will return the role with that name.
// Otherwise, it will return all roles.
func (svc *IAMService) GetRoles(ctx context.Context, input *ascTypes.GetRolesInput) ([]types.Role, error) {
	// TODO: Combine GetRole and ListRolePolicies into a single call.
	var roles []types.Role

	// Determine whether to call ListRoles or GetRole based on the input.
	if input.RoleName != "" {
		output, err := svc.Client.GetRole(ctx, &iam.GetRoleInput{
			RoleName: &input.RoleName,
		})
		if err != nil {
			return nil, err
		}
		roles = append(roles, *output.Role)
	} else {
		output, err := svc.Client.ListRoles(ctx, &iam.ListRolesInput{
			MaxItems:   &input.MaxItems,
			// Only set the path prefix if it is provided.
			PathPrefix: utils.StringPtr(input.PathPrefix),
		})
		if err != nil {
			return nil, err
		}

		roles = append(roles, output.Roles...)
	}

	return roles, nil
}

// GetRolePolicies gets the inline and managed policies for a role.
func (svc *IAMService) GetRoleInlinePolicies(ctx context.Context, input *ascTypes.GetRolePoliciesInput) ([]types.Policy, error) {
	output, err := svc.Client.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
		RoleName: &input.RoleName,
	})
	if err != nil {
		return nil, err
	}

	policies := make([]types.Policy, len(output.PolicyNames))
	for i, policyName := range output.PolicyNames {
		policies[i] = types.Policy{
			PolicyName: &policyName,
		}
	}

	return policies, nil
}

// GetRoleManagedPolicies gets the managed policies for a role.
func (svc *IAMService) GetRoleManagedPolicies(ctx context.Context, input *ascTypes.GetRoleManagedPoliciesInput) ([]types.Policy, error) {
	output, err := svc.Client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: &input.RoleName,
	})
	if err != nil {
		return nil, err
	}

	policies := make([]types.Policy, len(output.AttachedPolicies))
	for i, policy := range output.AttachedPolicies {
		policies[i] = types.Policy{
			PolicyName: policy.PolicyName,
		}
	}

	return policies, nil
}