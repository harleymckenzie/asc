"""
Test cases for the RDS module.

This module contains the test cases for the RDS module.
"""
import pytest
from unittest.mock import patch
from asc.services import rds
from .test_utils import setup_args


@patch('asc.services.rds.create_boto_session')
@pytest.mark.parametrize("displayed_tags, endpoint", [("Name", True), ("Name,Environment", False), (None, False)])
def test_list_rds_instances(create_boto_session, displayed_tags, endpoint, capsys):
    """
    Test case for the list_rds_instances function.

    This test verifies that the list_rds_instances function correctly retrieves
    and prints the details of RDS instances.

    Args:
        create_boto_session: Mocked function to create a Boto session.
        displayed_tags: Tags to display.
        capsys: Pytest fixture to capture stdout and stderr.

    Returns:
        None
    """
    args = setup_args(displayed_tags=displayed_tags)
    args.endpoint = endpoint

    mock_session = create_boto_session.return_value
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
                "DBInstanceClass": "db.t2.micro",
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
                "DBInstanceClass": "db.t2.micro",
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
                "ReaderEndpoint": "db-cluster-1.us-east-1.rds.amazonaws.com",
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

    if displayed_tags == "Name" and endpoint:
        assert out == (
            "Name                       Identifier                 Type         State      Endpoint                                       Role\n"
            "-------------------------  -------------------------  -----------  ---------  ---------------------------------------------  ------\n"
            "db-instance-aurora-reader  db-instance-aurora-reader  db.t2.micro  available  db-cluster-1.us-east-1.rds.amazonaws.com       Reader\n"
            "db-instance-aurora-writer  db-instance-aurora-writer  db.t2.micro  available  db-cluster-1.us-east-1.rds.amazonaws.com       Writer\n"
            "db-instance-mysql          db-instance-mysql          db.t2.micro  available  db-instance-mysql.us-east-1.rds.amazonaws.com\n"
        )
    elif displayed_tags == "Name,Environment" and not endpoint:
        assert out == (
            "Name                       Identifier                 Type         State      Role    Environment\n"
            "-------------------------  -------------------------  -----------  ---------  ------  --------------\n"
            "db-instance-aurora-reader  db-instance-aurora-reader  db.t2.micro  available  Reader  production\n"
            "db-instance-aurora-writer  db-instance-aurora-writer  db.t2.micro  available  Writer  production\n"
            "db-instance-mysql          db-instance-mysql          db.t2.micro  available          pre-production\n"
        )
    else:
        assert out == (
            "Identifier                 Type         State      Role\n"
            "-------------------------  -----------  ---------  ------\n"
            "db-instance-aurora-reader  db.t2.micro  available  Reader\n"
            "db-instance-aurora-writer  db.t2.micro  available  Writer\n"
            "db-instance-mysql          db.t2.micro  available\n"
            )
