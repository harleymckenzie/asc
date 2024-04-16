"""
Elasticache Redis service.

This module contains the functions for the Redis service.

Functions:
- add_subparsers(subparsers, global_parser) -> None
- list_redis_instances(args)
"""

from typing import Any
import logging
from ..common import (
    subparser_register,
    create_boto_session,
    print_as_table,
    apply_tags,
)


@subparser_register("redis")
def add_subparsers(subparsers: Any, global_parser: Any) -> None:
    """
    Add subparsers for common commands.

    Args:
        subparsers: The subparsers object from the main parser.
        global_parser: The global parser containing common arguments.
    """
    redis_parser = subparsers.add_parser(
        "redis",
        help="Redis service",
        description="Redis service",
        epilog="""Example: asc redis ls""",
        parents=[global_parser],
    )
    redis_parser.set_defaults(func=lambda args: redis_parser.print_help())
    redis_subparsers = redis_parser.add_subparsers(
        help="", metavar="subcommand", dest="subcommand"
    )

    # Redis list subcommand
    redis_list_parser = redis_subparsers.add_parser(
        "ls",
        help="List Redis instances",
        description="List Redis instances",
        epilog="""Example: asc redis ls""",
        parents=[global_parser],
    )
    redis_list_parser.add_argument(
        "--endpoint",
        "-e",
        help="Display endpoint in output.",
        action="store_true",
    )
    redis_list_parser.add_argument(
        "--sort-by",
        help="Sort the output by a specific key",
        default="Cluster Id",
    )
    redis_list_parser.add_argument(
        "--sort-order",
        help="Specify sort order: 'asc' for ascending or 'desc' for descending",
        default="asc",
    )
    redis_list_parser.set_defaults(func=list_redis_instances)


def list_redis_instances(args: Any) -> None:
    """
    Lists Redis instances based on given arguments.

    Args:
        args: The arguments received from the command-line input.

    Prints:
        A table displaying the details of all Redis instances.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    ec_client = session.client("elasticache")
    displayed_tags_list = args.config.get(
        "asc", "displayed_tags", fallback=""
    ).split(",")
    instance_list = []

    try:
        response = ec_client.describe_cache_clusters(ShowCacheNodeInfo=True)
    except Exception as e:
        logging.error("Failed to list Redis instances: %s", e)
        exit(1)

    # Cluster tags aren't returned in the response,
    # so we need to fetch them separately
    for cluster in response["CacheClusters"]:
        try:
            cluster["tag_response"] = ec_client.list_tags_for_resource(
                ResourceName=cluster["ARN"]
            )
        except ec_client.exceptions.CacheClusterNotFoundFault:
            # Catch CacheClusterNotFoundFault exceptions
            # if creation is in progress
            cluster["tag_response"] = {"TagList": []}
        except Exception as e:
            logging.error(
                "Error while listing tags for %s: %s",
                cluster["CacheClusterId"],
                e,
            )

        for instance_data in cluster["CacheNodes"]:
            instance_data["TagList"] = cluster["tag_response"]["TagList"]
            instance = {
                "Cluster Id": cluster["CacheClusterId"],
                "Type": cluster["CacheNodeType"],
                "Status": cluster["CacheClusterStatus"],
            }
            if args.endpoint:
                instance["Endpoint"] = instance_data["Endpoint"]["Address"]

            instance = apply_tags(instance, instance_data, displayed_tags_list)
            instance_list.append(instance)

    key_order = [
        "Cluster Id",
        "Type",
        "Status",
        "Endpoint",
    ] + displayed_tags_list
    print_as_table(
        instance_list,
        key_order=key_order,
        sort_key=args.sort_by,
        sort_order=args.sort_order,
    )
