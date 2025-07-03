package types

import (
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type RoleColumnDef struct {
	GetValue func(*types.Role) string
}

type GetRolesInput struct {
	RoleName   string
	MaxItems   int32
	PathPrefix string
}

type GetRolePoliciesInput struct {
	RoleName string
}

type GetRoleManagedPoliciesInput struct {
	RoleName string
}