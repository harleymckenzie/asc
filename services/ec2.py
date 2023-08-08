"""
EC2 Service
"""
import boto3
from .common import print_as_table


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
    """
    instance_list = []
    displayed_tags_list = args.config.get('asc', 'displayed_tags').split(',')
    ec2_client = args.session.client('ec2')
    response = ec2_client.describe_instances()

    for reservation in response["Reservations"]:
        for ec2_instance in reservation["Instances"]:
            instance = {"Public IP": ec2_instance.get("PublicIpAddress", ""),
                        "Id": ec2_instance["InstanceId"],
                        "Type": ec2_instance["InstanceType"],
                        "State": ec2_instance["State"]["Name"]}
            
            # Add tags to instance dict
            for tag in ec2_instance.get("Tags", []):
                if tag["Key"] in displayed_tags_list:
                    instance[tag["Key"]] = tag["Value"]

            instance_list.append(instance)

    print_as_table(instance_list)
