import boto3
from moto import mock_ec2
import pytest
from unittest.mock import patch, MagicMock
from services.ec2 import list_ec2_instances

@pytest.mark.parametrize(
    "displayed_tags, expected_tags",
    [
        ("Name", ["Name"]),
        ("Name,Environment", ["Name", "Environment"])
    ]
)
@mock_ec2
def test_list_ec2_instances(displayed_tags, expected_tags):
    # Create a mock EC2 instance
    ec2_client = boto3.client("ec2", region_name="eu-west-1")
    ec2_client.run_instances(ImageId="ami-12345678", MinCount=1, MaxCount=1, 
                             TagSpecifications=[{'ResourceType': 'instance', 'Tags': 
                                                 [{'Key': 'Name', 'Value': 'test1'}, 
                                                  {'Key': 'Owner', 'Value': 'UAT'},
                                                  {'Key': 'Environment', 'Value': 'Testing'}]}])
    
    # Fetch the created instances
    reservations = ec2_client.describe_instances()
    instance = reservations['Reservations'][0]['Instances'][0]

    instance_id = instance['InstanceId']
    instance_type = instance['InstanceType']
    public_ip = instance.get('PublicIpAddress', '')  # It might not always have a public IP

    # Filter out the tag values based on expected_tags
    tag_values = [tag['Value'] for tag in instance['Tags'] if tag['Key'] in expected_tags]

    # Constructing the expected output based on actual instance details
    expected_output = (
        "  Public IP       Id                   Type      State    {}\n"
        "--------------  -------------------  --------  -------  ------\n"
        "{}  {}  {}  running  {}"
    ).format(",".join(expected_tags), public_ip, instance_id, instance_type, "  ".join(tag_values))

    # Create an argparse.Namespace object to pass as an argument
    args = MagicMock()
    args.session = boto3.Session(region_name="eu-west-1")
    args.config.get.return_value = displayed_tags  # Mimic the behavior of ConfigParser's get() method for the given section and option

    # Print a newline directly to format pytest output
    print()

    # Capture the print output
    output_captured = []

    # Mock the print function and capture its arguments
    with patch("builtins.print", side_effect=lambda *args: output_captured.append(" ".join(map(str, args)))) as mocked_print:
        list_ec2_instances(args)

    output_string = "\n".join(output_captured)
    for tag in displayed_tags.split(','):
        assert tag in output_string

    # Print the captured output for manual verification
    print(output_string)