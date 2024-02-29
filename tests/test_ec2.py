"""
Test cases for the EC2 service

This module contains the test cases for the EC2 service.

"""
from unittest.mock import patch
import pytest
from core.services import ec2
from .test_utils import setup_args


@patch('core.services.ec2.create_boto_session')
@pytest.mark.parametrize("displayed_tags", [("Name"), ("Name,Environment"), (None)])
def test_list_ec2_instances(create_boto_session, displayed_tags, capsys):
    """
    Test case for the list_ec2_instances function.

    This test verifies that the list_ec2_instances function correctly retrieves
    and prints the details of EC2 instances.

    Args:
        mock_args: Mocked command-line arguments.
        create_boto_session: Mocked function to create a Boto session.
        capsys: Pytest fixture to capture stdout and stderr.

    Returns:
        None
    """
    args = setup_args(displayed_tags=displayed_tags)
    mock_session = create_boto_session.return_value
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
                    }
                ]
            }
        ]
    }
    ec2.list_ec2_instances(args)
    out, err = capsys.readouterr()
    
    print("\n" + out)

    if displayed_tags == "Name":
        assert out == (
            "Name        Public IP      Id                   Type      State\n"
            "----------  -------------  -------------------  --------  -------\n"
            "web-server  201.32.58.128  i-0f76a2e0e6b8d6e2c  t2.micro  running\n"
        )
    elif displayed_tags == "Name,Environment":
        assert out == (
            "Name        Public IP      Id                   Type      State    Environment\n"
            "----------  -------------  -------------------  --------  -------  -------------\n"
            "web-server  201.32.58.128  i-0f76a2e0e6b8d6e2c  t2.micro  running  production\n"
        )
    else:
        assert out == (
            "Public IP      Id                   Type      State\n"
            "-------------  -------------------  --------  -------\n"
            "201.32.58.128  i-0f76a2e0e6b8d6e2c  t2.micro  running\n"
        )
