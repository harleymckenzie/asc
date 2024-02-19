import pytest
from moto import mock_rds
import boto3
from asc.services import rds
from .test_utils import setup_args, run_and_capture_output

@pytest.fixture
def mock_rds_client():
    with mock_rds():
        yield boto3.client("rds", region_name="eu-west-1")

@pytest.mark.parametrize("displayed_tags", [("Name"), ("Name,Environment"), (None)])
def test_list_rds_instances(mock_rds_client, displayed_tags):
    # Create a mock RDS instance
    mock_rds_client.create_db_instance(
        DBInstanceIdentifier="test-instance",
        AllocatedStorage=20,
        DBInstanceClass="db.t2.micro",
        Engine="mysql",
        MasterUsername="admin",
        MasterUserPassword="password",
        Tags=[
            {"Key": "Name", "Value": "test-instance"},
            {"Key": "Environment", "Value": "Production"},
        ],
    )

    args = setup_args(displayed_tags)
    output = run_and_capture_output(rds.list_rds_instances, args)

    # Check for expected headers
    expected_headers = ["Identifier", "Type", "State"]
    assert all(header in output[0] for header in expected_headers)

    # Confirm that specified and unspecified tags are (not) in the output
    displayed_tags_set = set(displayed_tags.split(",")) if displayed_tags else set()
    unspecified_tags = {"Name", "Environment"} - displayed_tags_set
    for tag in displayed_tags_set:
        assert tag in output[0]
    for tag in unspecified_tags:
        assert tag not in output[0]
