"""
RDS service.

This module contains functions to interact with Amazon RDS service.

Functions:
- add_subparsers(subparsers, global_parser): Adds subparsers for RDS commands.
- list_rds_instances(args): Lists RDS instances.
"""
from ..common import (
    subparser_register,
    create_boto_session,
    print_as_table,
    apply_tags
)


@subparser_register('rds')
def add_subparsers(subparsers, global_parser):
    """
    Adds subparsers for RDS commands to the main parser.

    Args:
        subparsers: The subparsers object from the main parser.
        global_parser: The global parser containing common arguments.
    """
    rds_parser = subparsers.add_parser(
        "rds",
        help="RDS service",
        description="RDS service",
        epilog="""Example: asc rds ls""",
        parents=[global_parser],
    )
    rds_parser.set_defaults(func=lambda args: rds_parser.print_help())
    rds_subparsers = rds_parser.add_subparsers(
        help='',
        metavar='subcommand',
        dest='subcommand'
    )

    rds_list_parser = rds_subparsers.add_parser(
        "ls",
        help="List RDS instances",
        description="List RDS instances",
        epilog="""Example: asc rds ls""",
        parents=[global_parser],
    )
    rds_list_parser.add_argument(
        "--endpoint", "-e",
        help="Display endpoint in output.",
        action="store_true"
    )
    rds_list_parser.set_defaults(func=list_rds_instances)


def list_rds_instances(args):
    """
    Lists RDS instances based on given arguments.

    Args:
        args: The arguments received from the command-line input.

    Prints:
        A table displaying the details of all RDS instances.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    rds_client = session.client("rds")
    displayed_tags_list = args.config.get(
        "asc", "displayed_tags", fallback="").split(",")
    instance_list = []

    try:
        response = rds_client.describe_db_instances()
        cluster_response = rds_client.describe_db_clusters()
    except Exception as e:
        print(f"Failed to list RDS instances: {e}")
        exit(1)

    for instance_data in response["DBInstances"]:
        instance = {
            "Identifier": instance_data["DBInstanceIdentifier"],
            "Type": instance_data["DBInstanceClass"],
            "State": instance_data["DBInstanceStatus"],
        }

        if "aurora-mysql" in instance_data["Engine"]:
            instance = get_aurora_instance(
                instance, instance_data, cluster_response, args
            )
        else:
            if args.endpoint:
                instance["Endpoint"] = instance_data["Endpoint"]["Address"]

        instance = apply_tags(instance, instance_data, displayed_tags_list)
        instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i["Identifier"])
    print_as_table(instances)


def get_aurora_instance(instance, instance_data, cluster_response, args):
    """
    Get details of an Aurora instance.

    Args:
        instance: The instance dictionary to be updated.
        instance_data: The instance data received from RDS.
        cluster_response: The response from describe_db_clusters.
        args: The arguments received from the command-line input.

    Returns:
        The updated instance dictionary.
    """
    for cluster in cluster_response["DBClusters"]:
        for member in cluster["DBClusterMembers"]:
            if member["DBInstanceIdentifier"] == instance_data["DBInstanceIdentifier"]:
                # Only add Endpoint if args.endpoint is set
                if args.endpoint:
                    instance["Endpoint"] = (
                        cluster["Endpoint"]
                        if member["IsClusterWriter"]
                        else cluster["ReaderEndpoint"]
                    )
                instance["Role"] = (
                    "Writer" if member["IsClusterWriter"] else "Reader"
                )

    return instance
