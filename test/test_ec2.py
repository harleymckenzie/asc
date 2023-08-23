import boto3
from moto import mock_ec2
import pytest
from unittest.mock import patch, MagicMock
from services.ec2 import list_ec2_instances
import configparser

@pytest.mark.parametrize(
    "displayed_tags",
    [
        ("Name"),
        ("Name,Environment"),
        (None)
    ]
)
@mock_ec2
def test_list_ec2_instances(displayed_tags):
    # Create a mock EC2 instance
    ec2_client = boto3.client("ec2", region_name="eu-west-1")
    ec2_client.run_instances(ImageId="ami-12345678", MinCount=1, MaxCount=1, 
                             TagSpecifications=[{'ResourceType': 'instance', 'Tags': 
                                                 [{'Key': 'Name', 'Value': 'test1'}, 
                                                  {'Key': 'Owner', 'Value': 'UAT'},
                                                  {'Key': 'Environment', 'Value': 'Testing'}]}])

    # Create an argparse.Namespace object to pass as an argument
    args = MagicMock()
    args.session = boto3.Session(region_name="eu-west-1")
    args.config = configparser.RawConfigParser()
    args.config.add_section('asc')
    if displayed_tags:
        args.config.set('asc', 'displayed_tags', displayed_tags)

    # Capture the print output
    output_captured = []

    # Mock the print function and capture its arguments
    with patch("builtins.print", side_effect=lambda *args: output_captured.append(" ".join(map(str, args)))) as mocked_print:
        list_ec2_instances(args)

    output_string = "\n".join(output_captured)
    
    # Confirm that tags specified in displayed_tags is in the output
    if displayed_tags:
        for tag in displayed_tags.split(','):
            assert tag in output_string

    # Print the captured output for manual verification
    print('\n' + output_string)