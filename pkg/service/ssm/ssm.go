package ssm

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
	"github.com/jedib0t/go-pretty/v6/table"
)

type SSMClientAPI interface {
	DescribeParameters(ctx context.Context, params *ssm.DescribeParametersInput, optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type SSMService struct {
	Client SSMClientAPI
}

type columnDef struct {
	id    string
	title string
}

func getParameterName(p *types.ParameterMetadata) string {
	return aws.ToString(p.Name)
}

func getParameterType(p *types.ParameterMetadata) string {
	return string(p.Type)
}

func getParameterValue(p *types.Parameter) string {
	return aws.ToString(p.Value)
}

var availableColumns = []columnDef{
	{id: "name", title: "Name"},
	{id: "type", title: "Type"},
	{id: "value", title: "Value"},
}

func NewSSMService(ctx context.Context, profile string) (*SSMService, error) {
	var cfg aws.Config
	var err error

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)
	return &SSMService{Client: client}, nil
}

func (svc *SSMService) ListParameters(ctx context.Context, selectedColumns []string) error {

	params, err := svc.Client.DescribeParameters(ctx, &ssm.DescribeParametersInput{})
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headerRow := make(table.Row, 0)
	for _, colID := range selectedColumns {
		for _, col := range availableColumns {
			if col.id == colID {
				headerRow = append(headerRow, col.title)
				break
			}
		}
	}
	t.AppendHeader(headerRow)

	for _, param := range params.Parameters {
		row := make(table.Row, 0)
		for _, colID := range selectedColumns {
			for _, col := range availableColumns {
				if col.id == colID {
					row = append(row, col.getValue(&param))
					break
				}
			}
		}
		t.AppendRow(row)
	}

	tableformat.SetStyle(t, true, false, nil)
	fmt.Println("total", len(params.Parameters))
	t.Render()
	return nil
}

func (svc *SSMService) GetParameters(ctx context.Context, parameterNames []string) error {
	for _, name := range parameterNames {
		param, err := svc.Client.GetParameter(ctx, &ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("failed to get parameter %s: %w", name, err)
		}

		fmt.Printf("%s = %s\n", name, aws.ToString(param.Parameter.Value))
	}
	return nil
}
