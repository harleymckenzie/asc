package organizations

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	ascTypes "github.com/harleymckenzie/asc/internal/service/organizations/types"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

// OrganizationsClientAPI is the interface for the Organizations client.
type OrganizationsClientAPI interface {
	ListAccounts(context.Context, *organizations.ListAccountsInput, ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error)
	ListAccountsForParent(context.Context, *organizations.ListAccountsForParentInput, ...func(*organizations.Options)) (*organizations.ListAccountsForParentOutput, error)
	ListRoots(context.Context, *organizations.ListRootsInput, ...func(*organizations.Options)) (*organizations.ListRootsOutput, error)
	ListOrganizationalUnitsForParent(context.Context, *organizations.ListOrganizationalUnitsForParentInput, ...func(*organizations.Options)) (*organizations.ListOrganizationalUnitsForParentOutput, error)
	ListTagsForResource(context.Context, *organizations.ListTagsForResourceInput, ...func(*organizations.Options)) (*organizations.ListTagsForResourceOutput, error)
	DescribeOrganization(context.Context, *organizations.DescribeOrganizationInput, ...func(*organizations.Options)) (*organizations.DescribeOrganizationOutput, error)
	DescribeAccount(context.Context, *organizations.DescribeAccountInput, ...func(*organizations.Options)) (*organizations.DescribeAccountOutput, error)
	DescribeOrganizationalUnit(context.Context, *organizations.DescribeOrganizationalUnitInput, ...func(*organizations.Options)) (*organizations.DescribeOrganizationalUnitOutput, error)
}

// OrganizationsService is the service for the Organizations client.
type OrganizationsService struct {
	Client OrganizationsClientAPI
}

// NewOrganizationsService creates a new Organizations service.
func NewOrganizationsService(ctx context.Context, profile string, region string) (*OrganizationsService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	client := organizations.NewFromConfig(cfg.Config)
	return &OrganizationsService{Client: client}, nil
}

// GetAccounts returns all accounts in the organization.
func (svc *OrganizationsService) GetAccounts(ctx context.Context, input *ascTypes.GetAccountsInput) ([]types.Account, error) {
	var accounts []types.Account
	paginator := organizations.NewListAccountsPaginator(svc.Client, &organizations.ListAccountsInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, output.Accounts...)
	}
	return accounts, nil
}

// GetAccountsForParent returns all accounts for a given parent OU or root.
func (svc *OrganizationsService) GetAccountsForParent(ctx context.Context, parentID string) ([]types.Account, error) {
	var accounts []types.Account
	paginator := organizations.NewListAccountsForParentPaginator(svc.Client, &organizations.ListAccountsForParentInput{
		ParentId: &parentID,
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, output.Accounts...)
	}
	return accounts, nil
}

// GetRoots returns all roots in the organization.
func (svc *OrganizationsService) GetRoots(ctx context.Context) ([]types.Root, error) {
	output, err := svc.Client.ListRoots(ctx, &organizations.ListRootsInput{})
	if err != nil {
		return nil, err
	}
	return output.Roots, nil
}

// GetOUsForParent returns all OUs for a given parent OU or root.
func (svc *OrganizationsService) GetOUsForParent(ctx context.Context, parentID string) ([]types.OrganizationalUnit, error) {
	var ous []types.OrganizationalUnit
	paginator := organizations.NewListOrganizationalUnitsForParentPaginator(svc.Client, &organizations.ListOrganizationalUnitsForParentInput{
		ParentId: &parentID,
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		ous = append(ous, output.OrganizationalUnits...)
	}
	return ous, nil
}

// GetOrganization returns details about the organization.
func (svc *OrganizationsService) GetOrganization(ctx context.Context) (*types.Organization, error) {
	output, err := svc.Client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, err
	}
	return output.Organization, nil
}

// GetAccount returns details about a specific account.
func (svc *OrganizationsService) GetAccount(ctx context.Context, input *ascTypes.GetAccountInput) (*types.Account, error) {
	output, err := svc.Client.DescribeAccount(ctx, &organizations.DescribeAccountInput{
		AccountId: &input.AccountID,
	})
	if err != nil {
		return nil, err
	}
	return output.Account, nil
}

// GetOU returns details about a specific organizational unit.
func (svc *OrganizationsService) GetOU(ctx context.Context, input *ascTypes.GetOUInput) (*types.OrganizationalUnit, error) {
	output, err := svc.Client.DescribeOrganizationalUnit(ctx, &organizations.DescribeOrganizationalUnitInput{
		OrganizationalUnitId: &input.OUID,
	})
	if err != nil {
		return nil, err
	}
	return output.OrganizationalUnit, nil
}

// GetTagsForResource returns tags for a given resource.
func (svc *OrganizationsService) GetTagsForResource(ctx context.Context, resourceID string) ([]types.Tag, error) {
	output, err := svc.Client.ListTagsForResource(ctx, &organizations.ListTagsForResourceInput{
		ResourceId: &resourceID,
	})
	if err != nil {
		return nil, err
	}
	return output.Tags, nil
}

// OUTreeNode represents an OU or root in the organization tree, with its children and accounts.
type OUTreeNode struct {
	ID       string
	Name     string
	Children []OUTreeNode
	Accounts []types.Account
}

// BuildOUTree recursively builds the OU tree starting from a given parent ID.
func (svc *OrganizationsService) BuildOUTree(ctx context.Context, parentID string, parentName string) (*OUTreeNode, error) {
	node := &OUTreeNode{
		ID:   parentID,
		Name: parentName,
	}

	// Get accounts for this parent
	accounts, err := svc.GetAccountsForParent(ctx, parentID)
	if err != nil {
		return nil, err
	}
	node.Accounts = accounts

	// Get child OUs
	ous, err := svc.GetOUsForParent(ctx, parentID)
	if err != nil {
		return nil, err
	}

	for _, ou := range ous {
		child, err := svc.BuildOUTree(ctx, *ou.Id, *ou.Name)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, *child)
	}

	return node, nil
}

// GetAccountsWithOUs walks the OU tree and returns all accounts with their parent OU path.
func (svc *OrganizationsService) GetAccountsWithOUs(ctx context.Context) ([]ascTypes.AccountWithOU, error) {
	roots, err := svc.GetRoots(ctx)
	if err != nil {
		return nil, err
	}

	var result []ascTypes.AccountWithOU
	for _, root := range roots {
		tree, err := svc.BuildOUTree(ctx, *root.Id, *root.Name)
		if err != nil {
			return nil, err
		}
		collectAccounts(tree, "", &result)
	}
	return result, nil
}

// collectAccounts recursively collects accounts from the tree with their OU path.
func collectAccounts(node *OUTreeNode, parentPath string, result *[]ascTypes.AccountWithOU) {
	currentPath := node.Name
	if parentPath != "" {
		currentPath = parentPath + " / " + node.Name
	}

	for _, account := range node.Accounts {
		*result = append(*result, ascTypes.AccountWithOU{
			Account: account,
			OUName:  node.Name,
			OUPath:  currentPath,
		})
	}

	for _, child := range node.Children {
		collectAccounts(&child, currentPath, result)
	}
}
