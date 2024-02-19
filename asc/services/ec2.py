"""
EC2 service.

This module contains the functions for the EC2 service.

Functions:
- add_subparsers(subparsers, global_parser): Adds subparsers for EC2 commands.
- list_ec2_instances(args): Lists EC2 instances.
"""
from ..common import subparser_register, print_as_table, apply_tags


@subparser_register('ec2')
def add_subparsers(subparsers, global_parser) -> None:
    """
    Adds subparsers to the given subparsers object for the ec2 service.

    Args:
        subparsers: The subparsers object to add the subcommands to.
        global_parser: The global parser object to inherit options from.
    """
    ec2_parser = subparsers.add_parser(
        "ec2",
        help="EC2 service",
        description="EC2 service",
        epilog="""Example: asc ec2 ls""",

        parents=[global_parser],
    )
    ec2_parser.set_defaults(func=lambda args: ec2_parser.print_help())
    ec2_subparsers = ec2_parser.add_subparsers(
        help="Description:",
        dest="subcommand"
    )

    ec2_list_parser = ec2_subparsers.add_parser(
        "ls",
        help="List EC2 instances",
        description="List EC2 instances",
        epilog="""Example: asc ec2 ls""",
        parents=[global_parser],
    )
    ec2_list_parser.set_defaults(func=list_ec2_instances)


def list_ec2_instances(args):
    """
    List EC2 instances.

    Args:
        args: Arguments passed by the user, including configuration details
              and options such as displaying endpoints.

    Prints:
        A table displaying the details of all EC2 instances.
    """
    ec2_client = args.session.client("ec2")
    instance_list = []
    displayed_tags_list = args.config.get(
        "asc", "displayed_tags", fallback="").split(",")

    
    try:
        response = ec2_client.describe_instances()
    except Exception as e:
        print(f"Failed to list EC2 instances: {e}")
        exit(1)

    for reservation in response["Reservations"]:
        for instance_data in reservation["Instances"]:
            instance = {
                "Public IP": instance_data.get("PublicIpAddress", ""),
                "Id": instance_data["InstanceId"],
                "Type": instance_data["InstanceType"],
                "State": instance_data["State"]["Name"]
            }
            
            instance = apply_tags(instance, instance_data, displayed_tags_list)
            instance_list.append(instance)

    print_as_table(instance_list)
