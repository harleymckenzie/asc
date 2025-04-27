package rds

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/jedib0t/go-pretty/v6/table"
)

type mockRDSClient struct {
	describeInstancesOutput *rds.DescribeDBInstancesOutput
	describeClustersOutput  *rds.DescribeDBClustersOutput
	err                     error
}

func (m *mockRDSClient) DescribeDBInstances(
	_ context.Context,
	params *rds.DescribeDBInstancesInput,
	_ ...func(*rds.Options),
) (*rds.DescribeDBInstancesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeInstancesOutput, nil
}

func (m *mockRDSClient) DescribeDBClusters(
	_ context.Context,
	params *rds.DescribeDBClustersInput,
	_ ...func(*rds.Options),
) (*rds.DescribeDBClustersOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.describeClustersOutput, nil
}

func TestListInstances(t *testing.T) {
	testCases := []struct {
		name      string
		instances []types.DBInstance
		clusters  []types.DBCluster
		err       error
		wantErr   bool
	}{
		{
			name: "mixed instance types with Aurora and RDS replicas",
			instances: []types.DBInstance{
				{
					// Standalone RDS instance
					DBInstanceIdentifier: aws.String("postgres-1"),
					DBInstanceStatus:     aws.String("available"),
					Engine:               aws.String("postgres"),
					DBInstanceClass:      aws.String("db.t3.micro"),
				},
				{
					// RDS Primary with read replica
					DBInstanceIdentifier:             aws.String("mysql-primary"),
					DBInstanceStatus:                 aws.String("available"),
					Engine:                           aws.String("mysql"),
					DBInstanceClass:                  aws.String("db.t3.small"),
					ReadReplicaDBInstanceIdentifiers: []string{"mysql-replica"},
				},
				{
					// RDS Read replica
					DBInstanceIdentifier:                  aws.String("mysql-replica"),
					DBInstanceStatus:                      aws.String("available"),
					Engine:                                aws.String("mysql"),
					DBInstanceClass:                       aws.String("db.t3.small"),
					ReadReplicaSourceDBInstanceIdentifier: aws.String("mysql-primary"),
				},
				{
					// Aurora MySQL Writer instance
					DBInstanceIdentifier: aws.String("aurora-writer"),
					DBInstanceStatus:     aws.String("available"),
					Engine:               aws.String("aurora-mysql"),
					DBInstanceClass:      aws.String("db.r6g.large"),
					DBClusterIdentifier:  aws.String("aurora-cluster"),
				},
				{
					// Aurora MySQL Reader instance
					DBInstanceIdentifier: aws.String("aurora-reader"),
					DBInstanceStatus:     aws.String("available"),
					Engine:               aws.String("aurora-mysql"),
					DBInstanceClass:      aws.String("db.r6g.large"),
					DBClusterIdentifier:  aws.String("aurora-cluster"),
				},
			},
			clusters: []types.DBCluster{
				{
					DBClusterIdentifier: aws.String("aurora-cluster"),
					Engine:              aws.String("aurora-mysql"),
					DBClusterMembers: []types.DBClusterMember{
						{
							DBInstanceIdentifier: aws.String("aurora-writer"),
							IsClusterWriter:      aws.Bool(true),
						},
						{
							DBInstanceIdentifier: aws.String("aurora-reader"),
							IsClusterWriter:      aws.Bool(false),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:      "empty response",
			instances: []types.DBInstance{},
			clusters:  []types.DBCluster{},
			wantErr:   false,
		},
		{
			name:      "api error",
			instances: nil,
			clusters:  nil,
			err:       &types.DBInstanceNotFoundFault{Message: aws.String("DB instance not found")},
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockRDSClient{
				describeInstancesOutput: &rds.DescribeDBInstancesOutput{
					DBInstances: tc.instances,
				},
				describeClustersOutput: &rds.DescribeDBClustersOutput{
					DBClusters: tc.clusters,
				},
				err: tc.err,
			}

			svc := &RDSService{
				Client: mockClient,
			}

			instances, err := svc.GetInstances(context.Background())
			if (err != nil) != tc.wantErr {
				t.Errorf("ListInstances() error = %v, wantErr %v", err, tc.wantErr)
			}

			if len(instances) != len(tc.instances) {
				t.Errorf("ListInstances() returned %d instances, want %d", len(instances), len(tc.instances))
			}

			for i, instance := range instances {
				if instance.DBInstanceIdentifier != tc.instances[i].DBInstanceIdentifier {
					t.Errorf("ListInstances() returned instance %d with ID %s, want %s", i, *instance.DBInstanceIdentifier, *tc.instances[i].DBInstanceIdentifier)
				}
			}
		})
	}
}

func TestTableOutput(t *testing.T) {
	instances := []types.DBInstance{
		{
			DBInstanceIdentifier: aws.String("postgres-1"),
			DBInstanceStatus:     aws.String("available"),
			Engine:               aws.String("postgres"),
			EngineVersion:        aws.String("14.7"),
			DBInstanceClass:      aws.String("db.t3.micro"),
		},
		{
			DBInstanceIdentifier: aws.String("aurora-writer"),
			DBInstanceStatus:     aws.String("available"),
			Engine:               aws.String("aurora-mysql"),
			EngineVersion:        aws.String("8.0.23"),
			DBInstanceClass:      aws.String("db.r6g.large"),
			DBClusterIdentifier:  aws.String("aurora-cluster"),
		},
	}

	clusters := []types.DBCluster{
		{
			DBClusterIdentifier: aws.String("aurora-cluster"),
			Engine:              aws.String("aurora-mysql"),
			DBClusterMembers: []types.DBClusterMember{
				{
					DBInstanceIdentifier: aws.String("aurora-writer"),
					IsClusterWriter:      aws.Bool(true),
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
			name:            "full instance details",
			selectedColumns: []string{"Cluster Identifier", "Identifier", "Status", "Engine", "Engine Version", "Size", "Role"},
			wantHeaders:     table.Row{"Cluster Identifier", "Identifier", "Status", "Engine", "Engine Version", "Size", "Role"},
			wantRowCount:    2,
		},
		{
			name:            "minimal columns",
			selectedColumns: []string{"Identifier", "Status", "Role"},
			wantHeaders:     table.Row{"Identifier", "Status", "Role"},
			wantRowCount:    2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rdsTable := &RDSTable{
				Instances:       instances,
				Clusters:        clusters,
				SelectedColumns: tc.selectedColumns,
			}

			// Test Headers
			headers := rdsTable.Headers()
			if len(headers) != len(tc.wantHeaders) {
				t.Errorf("Headers() returned %d columns, want %d", len(headers), len(tc.wantHeaders))
			}

			// Test Rows
			rows := rdsTable.Rows()
			if len(rows) != tc.wantRowCount {
				t.Errorf("Rows() returned %d rows, want %d", len(rows), tc.wantRowCount)
			}

			// Print the actual table output
			tw := table.NewWriter()
			tw.AppendHeader(headers)
			tw.AppendRows(rows)
			tw.SetStyle(rdsTable.TableStyle())
			t.Logf("\nTable Output:\n%s", tw.Render())
		})
	}
}
