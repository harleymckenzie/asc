package types

import (
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

// GetAccountsInput is the input for the GetAccounts method.
type GetAccountsInput struct {
	ParentID string
}

// GetAccountInput is the input for the GetAccount method.
type GetAccountInput struct {
	AccountID string
}

// GetOUInput is the input for the GetOU method.
type GetOUInput struct {
	OUID string
}

// AccountWithOU represents an account with its parent OU path information.
type AccountWithOU struct {
	types.Account
	OUName string
	OUPath string
}
