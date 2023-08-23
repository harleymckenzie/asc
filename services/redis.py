from .common import print_as_table


def add_subparsers(subparsers, global_parser):
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
        help="Description:", dest="subcommand"
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
        "--endpoint", "-e",
        help="Display endpoint in output.",
        action="store_true"
    )
    redis_list_parser.set_defaults(func=list_redis_instances)


def list_redis_instances(args):
    """
    Lists Redis instances based on given arguments.

    Args:
        args: The arguments received from the command-line input.

    Prints:
        A table displaying the details of all Redis instances.
    """
    instance_list = []
    cluster_instance_tags = {}
    displayed_tags_list = args.config.get("asc", "displayed_tags").split(",")
    ec_client = args.session.client("elasticache")
    response = ec_client.describe_cache_clusters(ShowCacheNodeInfo=True)

    for cluster in response["CacheClusters"]:
        if displayed_tags_list:
            cluster_tags = ec_client.list_tags_for_resource(
                ResourceName=cluster["ARN"]
            )
            for tag in cluster_tags["TagList"]:
                # Store tags in cluster dict
                if tag["Key"] in displayed_tags_list:
                    cluster_instance_tags[tag["Key"]] = tag["Value"]

        for node in cluster["CacheNodes"]:
            instance = {
                "Name": cluster["CacheClusterId"],
                "Type": cluster["CacheNodeType"],
                "Status": cluster["CacheClusterStatus"],
            }
            instance.update(cluster_instance_tags)
            if args.endpoint:
                instance["Endpoint"] = node["Endpoint"]["Address"]

            instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i["Name"])
    print_as_table(instances)
