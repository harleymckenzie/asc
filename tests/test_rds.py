"""
Unit tests for the RDS service module.

Functions:
- test_list_rds_instances: Test the list_rds_instances function.
"""
import pytest
from asc.services import rds
from .test_utils import setup_args, run_and_capture_output, mock_rds_client


@pytest.mark.parametrize("displayed_tags", [("Name"), ("Name,Environment"), (None)])
def test_list_rds_instances(mock_rds_client, displayed_tags):
    """
    Test the list_rds_instances function.
    """
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

    # Print output for debugging
    print("\n" + output[0])

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