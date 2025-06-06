package types

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type ColumnDef struct {
	GetValue func(*types.Stack) string
}

type GetStacksInput struct {

	// The names of the stacks to get
	StackName *string
}
