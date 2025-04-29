package cloudformation

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/jedib0t/go-pretty/v6/table"
)

type mockCloudFormationClient struct {
	describeStacksOutput *cloudformation.DescribeStacksOutput
	err                  error
}

func (m *mockCloudFormationClient) DescribeStacks(
	_ context.Context,
	params *cloudformation.DescribeStacksInput,
	_ ...func(*cloudformation.Options),
) (*cloudformation.DescribeStacksOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeStacksOutput, nil
}

func TestListStacks(t *testing.T) {
	testCases := []struct {
		name    string
		stacks  []types.Stack
		err     error
		wantErr bool
	}{
		{
			name: "mixed stack states",
			stacks: []types.Stack{
				{
					StackName:   aws.String("test-stack-1"),
					StackStatus: types.StackStatusCreateComplete,
					StackId:     aws.String("arn:aws:cloudformation:region:account:stack/test-stack-1/id1"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Environment"),
							Value: aws.String("Production"),
						},
					},
				},
				{
					StackName:   aws.String("test-stack-2"),
					StackStatus: types.StackStatusUpdateInProgress,
					StackId:     aws.String("arn:aws:cloudformation:region:account:stack/test-stack-2/id2"),
					Tags: []types.Tag{
						{
							Key:   aws.String("Environment"),
							Value: aws.String("Staging"),
						},
					},
				},
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:    "empty response",
			stacks:  []types.Stack{},
			err:     nil,
			wantErr: false,
		},
		{
			name:    "api error",
			stacks:  nil,
			err:     errors.New("Stack does not exist"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockCloudFormationClient{
				describeStacksOutput: &cloudformation.DescribeStacksOutput{
					Stacks: tc.stacks,
				},
				err: tc.err,
			}

			svc := &CloudFormationService{
				Client: mockClient,
			}

			stacks, err := svc.GetStacks(context.Background())
			if (err != nil) != tc.wantErr {
				t.Errorf("ListStacks() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(stacks) != len(tc.stacks) {
				t.Errorf("ListStacks() returned %d stacks, want %d", len(stacks), len(tc.stacks))
			}

			for i, stack := range stacks {
				if stack.StackName != tc.stacks[i].StackName {
					t.Errorf("ListStacks() returned stack %d with name %s, want %s", i, *stack.StackName, *tc.stacks[i].StackName)
				}
			}
		})
	}
}

func TestTableOutput(t *testing.T) {
	stacks := []types.Stack{
		{
			StackName:   aws.String("test-stack-1"),
			StackStatus: types.StackStatusCreateComplete,
			StackId:     aws.String("arn:aws:cloudformation:region:account:stack/test-stack-1/id1"),
			Description: aws.String("Production infrastructure stack"),
			Tags: []types.Tag{
				{
					Key:   aws.String("Environment"),
					Value: aws.String("Production"),
				},
			},
		},
		{
			StackName:   aws.String("test-stack-2"),
			StackStatus: types.StackStatusUpdateInProgress,
			StackId:     aws.String("arn:aws:cloudformation:region:account:stack/test-stack-2/id2"),
			Description: aws.String("Staging infrastructure stack"),
			Tags: []types.Tag{
				{
					Key:   aws.String("Environment"),
					Value: aws.String("Staging"),
				},
			},
		},
	}

	testCases := []struct {
		name            string
		selectedColumns []string
		wantHeaders     table.Row
		wantRowCount    int
	}{
		{
			name:            "full stack details",
			selectedColumns: []string{"Stack Name", "Status", "Description"},
			wantHeaders:     table.Row{"Stack Name", "Status", "Description"},
			wantRowCount:    2,
		},
		{
			name:            "minimal columns",
			selectedColumns: []string{"Stack Name", "Status"},
			wantHeaders:     table.Row{"Stack Name", "Status"},
			wantRowCount:    2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfTable := &CloudFormationTable{
				Stacks:          stacks,
				SelectedColumns: tc.selectedColumns,
			}

			// Test Headers
			headers := cfTable.Headers()
			if len(headers) != len(tc.wantHeaders) {
				t.Errorf("Headers() returned %d columns, want %d", len(headers), len(tc.wantHeaders))
			}

			// Test Rows
			rows := cfTable.Rows()
			if len(rows) != tc.wantRowCount {
				t.Errorf("Rows() returned %d rows, want %d", len(rows), tc.wantRowCount)
			}

			// Print the actual table output
			tw := table.NewWriter()
			tw.AppendHeader(headers)
			tw.AppendRows(rows)
			tw.SetStyle(cfTable.TableStyle())
			t.Logf("\nTable Output:\n%s", tw.Render())
		})
	}
}
