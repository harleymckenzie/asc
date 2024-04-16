"""
Test cases for the EC2 service

This module contains the test cases for the EC2 service.
"""
from unittest.mock import patch
import logging
import pytest
from core.services.ec2 import add_subparsers, list_ec2_instances
from .test_utils import setup_parser, setup_config


@pytest.fixture
def mock_ec2_session(mocker):
    mock_session = mocker.patch('core.services.ec2.create_boto_session')
    mock_client = mock_session.return_value.client.return_value
    return mock_client


@pytest.mark.parametrize(
    "arg_list, displayed_tags, expected_output",
    [
        # Each tuple represents: argument list, tags to display, expected output
        (
            ['ec2', 'ls', '--profile', 'default', '--region', 'us-west-2'],
            "Name",  # Displayed tags
            (
                "Name        Id                   Type        State    Public IP\n"
                "----------  -------------------  ----------  -------  ---------------\n"
                "            i-2c76a2e0e6b8d6e2c  r7i.xlarge  running  112.112.113.115\n"
                "app-server  i-1c76a2e0e6b8d6e2c  r7i.xlarge  running  112.112.113.114\n"
                "web-server  i-0f76a2e0e6b8d6e2c  t2.micro    running  201.32.58.128"
            )
        ),
        (
            ['ec2', 'ls', '--profile', 'admin', '--region', 'eu-central-1', '--sort-by', 'Public IP'],
            "Name,Environment",  # Displayed tags
            (
                "Name        Id                   Type        State    Public IP        Environment\n"
                "----------  -------------------  ----------  -------  ---------------  -------------\n"
                "app-server  i-1c76a2e0e6b8d6e2c  r7i.xlarge  running  112.112.113.114\n"
                "            i-2c76a2e0e6b8d6e2c  r7i.xlarge  running  112.112.113.115  production\n"
                "web-server  i-0f76a2e0e6b8d6e2c  t2.micro    running  201.32.58.128    production"
            )
        ),
        (
            ['ec2', 'ls', '--profile', 'admin', '--region', 'eu-west-1', '--sort-by', 'Environment', '--sort-order', 'desc'],
            "Owner,Environment",  # Displayed tags
            (
                "Id                   Type        State    Public IP        Owner       Environment\n"
                "-------------------  ----------  -------  ---------------  ----------  -------------\n"
                "i-0f76a2e0e6b8d6e2c  t2.micro    running  201.32.58.128    John Doe    production\n"
                "i-2c76a2e0e6b8d6e2c  r7i.xlarge  running  112.112.113.115  John Doe    production\n"
                "i-1c76a2e0e6b8d6e2c  r7i.xlarge  running  112.112.113.114  Sarah Jane\n"
            )
        ),
    ],
    ids=[
        "No sorting or tags",
        "Sort by Public IP, Tags: Name,Environment",
        "Sort by Owner, Tags: Owner,Environment"
    ]
)
@patch('core.services.ec2.create_boto_session')
def test_list_ec2_instances(mock_create_boto_session, arg_list,
                            displayed_tags, expected_output, capsys):
    """
    Test case for the list_ec2_instances function.

    This test verifies that the list_ec2_instances function correctly retrieves
    and prints the details of EC2 instances.

    Args:
        create_boto_session: Mocked function to create a Boto session.
        capsys: Pytest fixture to capture stdout and stderr.

    Returns:
        None
    """
    args = setup_parser(add_subparsers, arg_list)
    args.config = setup_config(displayed_tags)

    mock_session = mock_create_boto_session.return_value
    mock_client = mock_session.client.return_value
    mock_client.describe_instances.return_value = {
        "Reservations": [
            {
                "Instances": [
                    {
                        "InstanceId": "i-0f76a2e0e6b8d6e2c",
                        "InstanceType": "t2.micro",
                        "State": {"Name": "running"},
                        "PublicIpAddress": "201.32.58.128",
                        "Tags": [
                            {"Key": "Name", "Value": "web-server"},
                            {"Key": "Environment", "Value": "production"},
                            {"Key": "Owner", "Value": "John Doe"}
                        ]
                    },
                    {
                        "InstanceId": "i-1c76a2e0e6b8d6e2c",
                        "InstanceType": "r7i.xlarge",
                        "State": {"Name": "running"},
                        "PublicIpAddress": "112.112.113.114",
                        "Tags": [
                            {"Key": "Name", "Value": "app-server"},
                            {"Key": "Owner", "Value": "Sarah Jane"}
                        ]
                    },
                    {
                        "InstanceId": "i-2c76a2e0e6b8d6e2c",
                        "InstanceType": "r7i.xlarge",
                        "State": {"Name": "running"},
                        "PublicIpAddress": "112.112.113.115",
                        "Tags": [
                            {"Key": "Environment", "Value": "production"},
                            {"Key": "Owner", "Value": "John Doe"}
                        ]
                    }
                ]
            }
        ]
    }
    list_ec2_instances(args)
    out, _ = capsys.readouterr()
    logging.info("\n" + out)
    print("\n" + out)
    assert out.strip() == expected_output.strip()

    # Verify that external calls were made as expected
    mock_create_boto_session.assert_called_once_with(profile=args.profile,
                                                     region=args.region)
    mock_client.describe_instances.assert_called_once()
