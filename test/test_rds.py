import boto3
from moto import mock_rds
import pytest
from unittest.mock import patch, MagicMock
from services.rds import list_rds_instances
import configparser


@pytest.mark.parametrize(
    "displayed_tags, display_endpoint",
    [("Name", True), ("Name,Environment", False), (None, True)],
)
@mock_rds
def test_list_rds_instances(displayed_tags, display_endpoint):
    # Create a mock RDS instance
    rds_client = boto3.client("rds", region_name="eu-west-1")
    rds_client.create_db_instance(
        DBInstanceIdentifier="test1",
        DBInstanceClass="db.t2.micro",
        Engine="mysql",
        Tags=[
            {"Key": "Name", "Value": "test-db1"},
            {"Key": "Owner", "Value": "UAT"},
            {"Key": "Environment", "Value": "Testing"},
        ],
    )

    # Create a mock Aurora cluster
    rds_client.create_db_cluster(
        DBClusterIdentifier="test-cluster1",
        Engine="aurora-mysql",
        MasterUsername="test",
        MasterUserPassword="test1234",
        Tags=[
            {"Key": "Name", "Value": "test-cluster1"},
            {"Key": "Owner", "Value": "UAT"},
            {"Key": "Environment", "Value": "Testing"},
        ],
    )

    # Create a mock Aurora writer and reader instance
    rds_client.create_db_instance(
        DBInstanceIdentifier="test-aurora-instance1",
        DBClusterIdentifier="test-cluster1",
        DBInstanceClass="db.t2.micro",
        Engine="aurora-mysql",
        Tags=[
            {"Key": "Name", "Value": "test-cluster1-instance1"},
            {"Key": "Owner", "Value": "UAT"},
            {"Key": "Environment", "Value": "Testing"},
        ],
    )

    # Create a mock Aurora replica instance
    rds_client.create_db_instance_read_replica(
        DBInstanceIdentifier="test-replica1",
        SourceDBInstanceIdentifier="test-aurora-instance1",
        Tags=[
            {"Key": "Name", "Value": "test-cluster1-replica1"},
            {"Key": "Owner", "Value": "UAT"},
            {"Key": "Environment", "Value": "Testing"},
        ],
    )

    # Create an argparse.Namespace object to pass as an argument
    args = MagicMock()
    args.session = boto3.Session(region_name="eu-west-1")
    args.config = configparser.RawConfigParser()
    args.config.add_section("asc")
    if displayed_tags:
        args.config.set("asc", "displayed_tags", displayed_tags)
    args.endpoint = display_endpoint

    # Capture the print output
    output_captured = []

    # Mock the print function and capture its arguments
    with patch(
        "builtins.print",
        side_effect=lambda *args: output_captured.append(" ".join(map(str, args))),
    ) as mocked_print:
        list_rds_instances(args)

    output_string = "\n".join(output_captured)

    if displayed_tags:
        for tag in displayed_tags.split(","):
            assert tag in output_string

    # Print the captured output for manual verification
    print("\n" + output_string)
