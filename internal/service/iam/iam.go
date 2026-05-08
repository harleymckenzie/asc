package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

type IAMClientAPI interface {
	ListRoles(context.Context, *iam.ListRolesInput, ...func(*iam.Options)) (*iam.ListRolesOutput, error)
	GetRole(context.Context, *iam.GetRoleInput, ...func(*iam.Options)) (*iam.GetRoleOutput, error)
	ListPolicies(context.Context, *iam.ListPoliciesInput, ...func(*iam.Options)) (*iam.ListPoliciesOutput, error)
	GetPolicy(context.Context, *iam.GetPolicyInput, ...func(*iam.Options)) (*iam.GetPolicyOutput, error)
	ListUsers(context.Context, *iam.ListUsersInput, ...func(*iam.Options)) (*iam.ListUsersOutput, error)
	GetUser(context.Context, *iam.GetUserInput, ...func(*iam.Options)) (*iam.GetUserOutput, error)
	ListAttachedRolePolicies(context.Context, *iam.ListAttachedRolePoliciesInput, ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
	ListAttachedUserPolicies(context.Context, *iam.ListAttachedUserPoliciesInput, ...func(*iam.Options)) (*iam.ListAttachedUserPoliciesOutput, error)
	ListGroupsForUser(context.Context, *iam.ListGroupsForUserInput, ...func(*iam.Options)) (*iam.ListGroupsForUserOutput, error)
}

type IAMService struct {
	Client IAMClientAPI
}

func NewIAMService(ctx context.Context, profile string, region string) (*IAMService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	client := iam.NewFromConfig(cfg.Config)
	return &IAMService{Client: client}, nil
}

func (svc *IAMService) GetRoles(ctx context.Context) ([]types.Role, error) {
	var roles []types.Role
	input := &iam.ListRolesInput{}

	for {
		output, err := svc.Client.ListRoles(ctx, input)
		if err != nil {
			return nil, err
		}
		roles = append(roles, output.Roles...)
		if !output.IsTruncated {
			break
		}
		input.Marker = output.Marker
	}

	return roles, nil
}

func (svc *IAMService) GetRole(ctx context.Context, roleName string) (*types.Role, error) {
	output, err := svc.Client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return nil, err
	}
	return output.Role, nil
}

func (svc *IAMService) GetPolicies(ctx context.Context, includeAWSManaged bool) ([]types.Policy, error) {
	var policies []types.Policy
	input := &iam.ListPoliciesInput{
		Scope: "Local",
	}
	if includeAWSManaged {
		input.Scope = "All"
	}

	for {
		output, err := svc.Client.ListPolicies(ctx, input)
		if err != nil {
			return nil, err
		}
		policies = append(policies, output.Policies...)
		if !output.IsTruncated {
			break
		}
		input.Marker = output.Marker
	}

	return policies, nil
}

func (svc *IAMService) GetPolicy(ctx context.Context, identifier string) (*types.Policy, error) {
	if strings.HasPrefix(identifier, "arn:") {
		return svc.getPolicyByARN(ctx, identifier)
	}
	return svc.getPolicyByName(ctx, identifier)
}

func (svc *IAMService) getPolicyByARN(ctx context.Context, arn string) (*types.Policy, error) {
	output, err := svc.Client.GetPolicy(ctx, &iam.GetPolicyInput{
		PolicyArn: aws.String(arn),
	})
	if err != nil {
		return nil, err
	}
	return output.Policy, nil
}

func (svc *IAMService) getPolicyByName(ctx context.Context, name string) (*types.Policy, error) {
	policies, err := svc.GetPolicies(ctx, true)
	if err != nil {
		return nil, err
	}

	var matches []types.Policy
	for _, p := range policies {
		if aws.ToString(p.PolicyName) == name {
			matches = append(matches, p)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no policy found with name %q", name)
	case 1:
		return &matches[0], nil
	default:
		return nil, fmt.Errorf("multiple policies found with name %q, use an ARN instead", name)
	}
}

func (svc *IAMService) GetUsers(ctx context.Context) ([]types.User, error) {
	var users []types.User
	input := &iam.ListUsersInput{}

	for {
		output, err := svc.Client.ListUsers(ctx, input)
		if err != nil {
			return nil, err
		}
		users = append(users, output.Users...)
		if !output.IsTruncated {
			break
		}
		input.Marker = output.Marker
	}

	return users, nil
}

func (svc *IAMService) GetUser(ctx context.Context, userName string) (*types.User, error) {
	output, err := svc.Client.GetUser(ctx, &iam.GetUserInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		return nil, err
	}
	return output.User, nil
}

func (svc *IAMService) GetAttachedRolePolicies(ctx context.Context, roleName string) ([]types.AttachedPolicy, error) {
	var policies []types.AttachedPolicy
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}

	for {
		output, err := svc.Client.ListAttachedRolePolicies(ctx, input)
		if err != nil {
			return nil, err
		}
		policies = append(policies, output.AttachedPolicies...)
		if !output.IsTruncated {
			break
		}
		input.Marker = output.Marker
	}

	return policies, nil
}

func (svc *IAMService) GetAttachedUserPolicies(ctx context.Context, userName string) ([]types.AttachedPolicy, error) {
	var policies []types.AttachedPolicy
	input := &iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(userName),
	}

	for {
		output, err := svc.Client.ListAttachedUserPolicies(ctx, input)
		if err != nil {
			return nil, err
		}
		policies = append(policies, output.AttachedPolicies...)
		if !output.IsTruncated {
			break
		}
		input.Marker = output.Marker
	}

	return policies, nil
}

func (svc *IAMService) GetGroupsForUser(ctx context.Context, userName string) ([]types.Group, error) {
	var groups []types.Group
	input := &iam.ListGroupsForUserInput{
		UserName: aws.String(userName),
	}

	for {
		output, err := svc.Client.ListGroupsForUser(ctx, input)
		if err != nil {
			return nil, err
		}
		groups = append(groups, output.Groups...)
		if !output.IsTruncated {
			break
		}
		input.Marker = output.Marker
	}

	return groups, nil
}
