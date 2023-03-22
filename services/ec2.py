'''
EC2 Service
'''
import boto3
from .common import print_as_table, load_config

def add_subparsers(subparsers):
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
    """
    instance_list = []
    config = load_config()
    env_tag_key = config.get('asc', 'env_tag_key', fallback='Environment')
    ec2 = boto3.client('ec2')
    response = ec2.describe_instances()

    for reservation in response["Reservations"]:
        for ec2 in reservation["Instances"]:

            # If the instance has a name or environment tag set in config, use it
            if "Tags" in ec2:
                for tag in ec2["Tags"]:
                    if tag["Key"] == "Name":
                        instance_name = tag["Value"]
                    if tag["Key"] == env_tag_key:
                        instance_stack = tag["Value"]

            instance = {}
            instance["Name"] = instance_name if "instance_name" in locals() else ""
            if "instance_stack" in locals():
                instance["Environment"] = instance_stack
            instance["Public IP"] = ec2["PublicIpAddress"] if "PublicIpAddress" in ec2 else ""
            instance["Id"] = ec2["InstanceId"]
            instance["Type"] = ec2["InstanceType"]
            instance["State"] = ec2["State"]["Name"]
            instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i['Name'])
    print_as_table(instances)
