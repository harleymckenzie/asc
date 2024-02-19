"""
Unit tests for the EC2 service module.

Functions:
- test_list_ec2_instances: Test the list_ec2_instances function.
"""
import pytest
from moto import mock_ec2
import boto3
from asc.services import ec2
from .test_utils import setup_args, run_and_capture_output

@pytest.fixture
def mock_ec2_client():
    """
    Create a mock EC2 client.
    """
    with mock_ec2():
        yield boto3.client("ec2", region_name="eu-west-1")

@pytest.mark.parametrize("displayed_tags", [("Name"), ("Name,Environment"), (None)])
def test_list_ec2_instances(displayed_tags):
    """
    Test the list_ec2_instances function.
    """
    image_response = mock_ec2_client.describe_images()
    image_id = image_response['Images'][0]['ImageId']
    mock_ec2_client.run_instances(
        ImageId=image_id,
        MinCount=1,
        MaxCount=1,
        TagSpecifications=[
            {
                "ResourceType": "instance",
                "Tags": [
                    {"Key": "Name", "Value": "test-instance"},
                    {"Key": "Environment", "Value": "Production"},
                ],
            }
        ],
    )

    args = setup_args(displayed_tags)
    output = run_and_capture_output(ec2.list_ec2_instances, args)

    # Check for expected headers
    expected_headers = ["Public IP", "Id", "Type", "State"]
    assert all(header in output[0] for header in expected_headers)

    # Confirm that specified and unspecified tags are (not) in the output
    displayed_tags_set = set(displayed_tags.split(",")) if displayed_tags else set()
    unspecified_tags = {"Name", "Environment"} - displayed_tags_set
    for tag in displayed_tags_set:
        assert tag in output[0]
    for tag in unspecified_tags:
        assert tag not in output[0]
