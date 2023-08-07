"""
EC2 Service

This module provides functionality for interacting with Amazon Elastic Compute Cloud (EC2) service.
"""

import boto3
from .common import print_as_table, load_config


def add_subparsers(subparsers):
    """
    Add subparsers for EC2 commands
    """
    ec2_parser = subparsers.add_parser('ec2', help='EC2 service', description='EC2 service',
                                       epilog='''Example: asc ec2 ls''')
    ec2_parser.set_defaults(func=lambda args: ec2_parser.print_help())
    ec2_subparsers = ec2_parser.add_subparsers(help='Description:', dest='subcommand')

    ec2_list_parser = ec2_subparsers.add_parser('ls', help='List EC2 instances',
                                                description='List EC2 instances', epilog='''Example: asc ec2 ls''')
    ec2_list_parser.set_defaults(func=list_ec2_instances)


def list_ec2_instances(args):
    """
    List EC2 instances

    This function retrieves a list of EC2 instances and prints them as a table.

    Args:
        args: The arguments passed to the command.

    Returns:
        None
    """
    instance_list = []
    config = load_config()
    env_tag_key = config.get('asc', 'env_tag_key', fallback='Environment')
    ec2 = boto3.client('ec2')
    response = ec2.describe_instances()

    for reservation in response["Reservations"]:
        for ec2 in reservation["Instances"]:
            instance = {"Public IP": ec2["PublicIpAddress"] if "PublicIpAddress" in ec2 else "",
                        "Id": ec2["InstanceId"], "Type": ec2["InstanceType"], "State": ec2["State"]["Name"]}

            # If the instance has a name or environment tag set in config, use it
            if "Tags" in ec2:
                for tag in ec2["Tags"]:
                    if tag["Key"] == "Name":
                        instance = {"Name": tag["Value"], **instance}
                    if tag["Key"] == env_tag_key:
                        instance["Environment"] = tag["Value"]

            instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i['Name'])
    print_as_table(instances)
