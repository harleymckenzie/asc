"""
Test cases for the RDS module.

This module contains the test cases for the RDS module.
"""
from unittest.mock import patch
import pytest
from core.services import rds
from .test_utils import setup_parser, setup_config


@pytest.mark.parametrize(
    "arg_list, displayed_tags, expected_output",
    [
        (
            ['rds', 'ls', '--profile', 'test-profile', '--region', 'eu-west-1', '--endpoint'],
            "Name",
            "Name                       Identifier                 Type            State      Endpoint                                       Role\n"
            "-------------------------  -------------------------  --------------  ---------  ---------------------------------------------  ------\n"
            "db-instance-aurora-reader  db-instance-aurora-reader  db.r7g.xlarge   available  db-cluster-1-ro.us-east-1.rds.amazonaws.com    Reader\n"
            "db-instance-aurora-writer  db-instance-aurora-writer  db.t2.micro     available  db-cluster-1.us-east-1.rds.amazonaws.com       Writer\n"
            "db-instance-mysql          db-instance-mysql          db.r7g.4xlarge  available  db-instance-mysql.us-east-1.rds.amazonaws.com\n"
        ),
        (
            ['rds', 'ls', '--profile', 'test-profile', '--region', 'eu-west-1', '--sort-by', 'Identifier', '--sort-order', 'desc'],
            "Environment",
            "Identifier                 Type            State      Role    Environment\n"
            "-------------------------  --------------  ---------  ------  --------------\n"
            "db-instance-mysql          db.r7g.4xlarge  available          pre-production\n"
            "db-instance-aurora-writer  db.t2.micro     available  Writer  production\n"
            "db-instance-aurora-reader  db.r7g.xlarge   available  Reader  production\n"
        ),
        (
            ['rds', 'ls', '--profile', 'test-profile', '--region', 'eu-west-1', '--sort-by', 'Type', '--sort-order', 'desc'],
            "Name,Environment",
            "Name                       Identifier                 Type            State      Role    Environment\n"
            "-------------------------  -------------------------  --------------  ---------  ------  --------------\n"
            "db-instance-aurora-writer  db-instance-aurora-writer  db.t2.micro     available  Writer  production\n"
            "db-instance-aurora-reader  db-instance-aurora-reader  db.r7g.xlarge   available  Reader  production\n"
            "db-instance-mysql          db-instance-mysql          db.r7g.4xlarge  available          pre-production\n"
        )
    ],
    ids=[
        "Output endpoint, tags: Name",
        "tags: Environment",
        "Sort by Type, tags: Name,Environment"
    ]
)
@patch('core.services.rds.create_boto_session')
def test_list_rds_instances(mock_create_boto_session, arg_list,
                            displayed_tags, expected_output, capsys):
    """
    Test case for the list_rds_instances function.

    This test verifies that the list_rds_instances function correctly retrieves
    and prints the details of RDS instances.

    Args:
        mock_create_boto_session: Mocked function to create a Boto session.
        displayed_tags: Tags to display.
        capsys: Pytest fixture to capture stdout and stderr.

    Returns:
        None
    """
    args = setup_parser(rds.add_subparsers, arg_list)
    args.config = setup_config(displayed_tags)

    mock_session = mock_create_boto_session.return_value
    mock_client = mock_session.client.return_value
    mock_client.describe_db_instances.return_value = {
        "DBInstances": [
            {
                "DBInstanceIdentifier": "db-instance-aurora-writer",
                "DBInstanceClass": "db.t2.micro",
                "Engine": "aurora-mysql",
                "EngineVersion": "5.7.mysql_aurora.2.03.2",
                "DBClusterIdentifier": "db-cluster-1",
                "DBInstanceStatus": "available",
                "PubliclyAccessible": False,
                "Endpoint": {
                    "Address": "db-instance-aurora-writer.cluster-1.us-east-1.rds.amazonaws.com",
                    "Port": 3306,
                    "HostedZoneId": "Z1H1FL5HABSF5"
                },
                "TagList": [
                    {"Key": "Name", "Value": "db-instance-aurora-writer"},
                    {"Key": "Environment", "Value": "production"},
                    {"Key": "Owner", "Value": "John Doe"}
                ]
            },
            {
                "DBInstanceIdentifier": "db-instance-aurora-reader",
                "DBInstanceClass": "db.r7g.xlarge",
                "Engine": "aurora-mysql",
                "EngineVersion": "5.7.mysql_aurora.2.03.2",
                "DBClusterIdentifier": "db-cluster-1",
                "DBInstanceStatus": "available",
                "PubliclyAccessible": False,
                "Endpoint": {
                    "Address": "db-instance-aurora-reader.cluster-1.us-east-1.rds.amazonaws.com",
                    "Port": 3306,
                    "HostedZoneId": "Z1H1FL5HABSF5"
                },
                "TagList": [
                    {"Key": "Name", "Value": "db-instance-aurora-reader"},
                    {"Key": "Environment", "Value": "production"},
                    {"Key": "Owner", "Value": "John Doe"}
                ]
            },
            {
                "DBInstanceIdentifier": "db-instance-mysql",
                "DBInstanceClass": "db.r7g.4xlarge",
                "Engine": "mysql",
                "EngineVersion": "5.7.30",
                "DBInstanceStatus": "available",
                "PubliclyAccessible": False,
                "Endpoint": {
                    "Address": "db-instance-mysql.us-east-1.rds.amazonaws.com",
                    "Port": 3306,
                    "HostedZoneId": "Z1H1FL5HABSF5"
                },
                "TagList": [
                    {"Key": "Name", "Value": "db-instance-mysql"},
                    {"Key": "Environment", "Value": "pre-production"},
                    {"Key": "Owner", "Value": "John Doe"}
                ]
            }
        ]
    }
    mock_client.describe_db_clusters.return_value = {
        "DBClusters": [
            {
                "DBClusterIdentifier": "db-cluster-1",
                "DBClusterMembers": [{
                    "DBInstanceIdentifier": "db-instance-aurora-writer",
                    "IsClusterWriter": True
                }, {
                    "DBInstanceIdentifier": "db-instance-aurora-reader",
                    "IsClusterWriter": False
                }],
                "Endpoint": "db-cluster-1.us-east-1.rds.amazonaws.com",
                "Engine": "aurora-mysql",
                "Status": "available",
                "EngineVersion": "5.7.mysql_aurora.2.03.2",
                "ReaderEndpoint": "db-cluster-1-ro.us-east-1.rds.amazonaws.com",
                "TagList": [
                    {"Key": "Name", "Value": "db-cluster-1"},
                    {"Key": "Environment", "Value": "production"},
                    {"Key": "Owner", "Value": "John Doe"}
                ]
            }
        ]
    }
    rds.list_rds_instances(args)
    out, _ = capsys.readouterr()
    print("\n" + out)
    assert out.strip() == expected_output.strip()

    mock_create_boto_session.assert_called_once_with(
        profile=args.profile, region=args.region)
    mock_client.describe_db_instances.assert_called_once()
    mock_client.describe_db_clusters.assert_called_once()
