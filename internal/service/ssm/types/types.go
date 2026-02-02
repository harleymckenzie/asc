package types

// GetParameterInput defines input for GetParameter operation.
type GetParameterInput struct {
	Name    string
	Decrypt bool
}

// GetParametersInput defines input for GetParameters operation.
type GetParametersInput struct {
	Names   []string
	Decrypt bool
}

// GetParametersByPathInput defines input for GetParametersByPath operation.
type GetParametersByPathInput struct {
	Path      string
	Recursive bool
	Decrypt   bool
}

// PutParameterInput defines input for PutParameter operation.
type PutParameterInput struct {
	Name        string
	Value       string
	Type        string // String, StringList, SecureString
	Description string
	Overwrite   bool
	Tags        map[string]string
}

// DeleteParameterInput defines input for DeleteParameter operation.
type DeleteParameterInput struct {
	Name string
}

// DeleteParametersInput defines input for DeleteParameters operation.
type DeleteParametersInput struct {
	Names []string
}

// CopyParameterInput defines input for CopyParameter operation.
type CopyParameterInput struct {
	Source    string
	Dest      string
	Overwrite bool
}

// MoveParameterInput defines input for MoveParameter operation.
type MoveParameterInput struct {
	Source string
	Dest   string
}
