package cloudformation

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	ascTypes "github.com/harleymckenzie/asc/pkg/service/cloudformation/types"
	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

type CloudFormationTable struct {
	Stacks          []types.Stack
	SelectedColumns []string
	SortBy          string
}

type CloudFormationClientAPI interface {
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
}

type CloudFormationService struct {
	Client CloudFormationClientAPI
}

func availableColumns() map[string]ascTypes.ColumnDef {
	return map[string]ascTypes.ColumnDef{
		"Stack Name": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.StackName)
			},
		},
		"Status": {
			GetValue: func(i *types.Stack) string {
				return tableformat.ResourceState(string(i.StackStatus))
			},
		},
		"Description": {
			GetValue: func(i *types.Stack) string {
				return aws.ToString(i.Description)
			},
		},
	}
}

//
// Table functions
//

func (et *CloudFormationTable) Headers() table.Row {
	return tableformat.BuildHeaders(et.SelectedColumns)
}

func (et *CloudFormationTable) Rows() []table.Row {
	rows := []table.Row{}
	for _, stack := range et.Stacks {
		row := table.Row{}
		for _, colID := range et.SelectedColumns {
			row = append(row, availableColumns()[colID].GetValue(&stack))
		}
		rows = append(rows, row)
	}
	return rows
}

func (et *CloudFormationTable) ColumnConfigs() []table.ColumnConfig {
	return []table.ColumnConfig{
		{Name: "Stack Name", WidthMax: 80},
		// {Name: "Status", WidthMax: 15},
	}
}

func (et *CloudFormationTable) TableStyle() table.Style {
	return table.StyleRounded
}

//
// Service functions
//

func NewCloudFormationService(ctx context.Context, profile string, region string) (*CloudFormationService, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, err
	}

	return &CloudFormationService{
		Client: cloudformation.NewFromConfig(cfg),
	}, nil
}

func (svc *CloudFormationService) GetStacks(ctx context.Context) ([]types.Stack, error) {
	input := &cloudformation.DescribeStacksInput{}

	output, err := svc.Client.DescribeStacks(ctx, input)
	if err != nil {
		return nil, err
	}

	return output.Stacks, nil
}
