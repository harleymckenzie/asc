"""
Redis service
"""
import boto3
from .common import print_as_table, load_config


def add_subparsers(subparsers):
    """
    Add subparsers for Redis commands
    """
    redis_parser = subparsers.add_parser('redis', help='Redis service', description='Redis service',
                                         epilog='''Example: asc redis ls''')
    redis_parser.set_defaults(func=lambda args: redis_parser.print_help())
    redis_subparsers = redis_parser.add_subparsers(help='Description:', dest='subcommand')

    # Redis list subcommand
    redis_list_parser = redis_subparsers.add_parser('ls', help='List Redis instances',
                                                    description='List Redis instances',
                                                    epilog='''Example: asc redis ls''')
    redis_list_parser.add_argument('--endpoint', '-e', help='Display endpoint in output.', action='store_true')
    redis_list_parser.set_defaults(func=list_redis_instances)

    # Redis tag subcommand
    redis_tag_parser = redis_subparsers.add_parser('tag')
    redis_tag_parser.add_argument('name', help='Redis instance name')
    redis_tag_parser.add_argument('environment', help='Environment')
    redis_tag_parser.set_defaults(func=tag_redis)


def list_redis_instances(args):
    """
    List Redis instances
    """
    instance_list = []
    displayed_tags_list = args.config.get('asc', 'displayed_tags').split(',')
    elasticache_client = args.session.client('elasticache')
    response = elasticache_client.describe_cache_clusters(ShowCacheNodeInfo=True)

    # Loop through clusters and nodes
    for cluster in response["CacheClusters"]:
        # If tags are present in displayed_tags_list, retrieve them from the cluster
        if displayed_tags_list:
            cluster_tags = elasticache_client.list_tags_for_resource(ResourceName=cluster["ARN"])
            cluster_instance_tags = {}
            for tag in cluster_tags["TagList"]:
                # Store tags in cluster dict
                if tag["Key"] in displayed_tags_list:
                    cluster_instance_tags[tag["Key"]] = tag["Value"]

        # Loop through nodes in cluster, as there can be multiple nodes per cluster
        # If there is a cluster dict, add tags to instance dict
        for node in cluster["CacheNodes"]:
            instance = {"Name": cluster["CacheClusterId"],
                        "Type": cluster["CacheNodeType"],
                        "Status": cluster["CacheClusterStatus"]}
            instance.update(cluster_instance_tags)
            if args.endpoint:
                instance["Endpoint"] = node["Endpoint"]["Address"]

            instance_list.append(instance)


    instances = sorted(instance_list, key=lambda i: i['Name'])
    print_as_table(instances)


def tag_redis(args):
    """
    Tag Redis instance
    """
    elasticache = boto3.client('elasticache')
