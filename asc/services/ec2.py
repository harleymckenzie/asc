"""
EC2 service.

This module contains the functions for the EC2 service.

Functions:
- add_subparsers(subparsers, global_parser): Adds subparsers for EC2 commands.
- list_ec2_instances(args): Lists EC2 instances.
"""
from ..common import subparser_register, print_as_table


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
    instance_list = []

    # Store tags to display in the output if they've been set in the config
    if "displayed_tags" in args.config["asc"]:
        displayed_tags_list = args.config["asc"]["displayed_tags"].split(",")
    # Set an empty list if the config hasn't been set
    else:
        displayed_tags_list = []

    ec2_client = args.session.client("ec2")
    response = ec2_client.describe_instances()

    for reservation in response["Reservations"]:
        for ec2_instance in reservation["Instances"]:
            instance = {
                "Public IP": ec2_instance.get("PublicIpAddress", ""),
                "Id": ec2_instance["InstanceId"],
                "Type": ec2_instance["InstanceType"],
                "State": ec2_instance["State"]["Name"],
            }

            # Add tags to instance dict
            for tag in ec2_instance.get("Tags", []):
                if tag["Key"] in displayed_tags_list:
                    if "Name" == tag["Key"]:
                        instance = {tag["Key"]: tag["Value"], **instance}
                    else:
                        instance[tag["Key"]] = tag["Value"]

            instance_list.append(instance)

    print_as_table(instance_list)
