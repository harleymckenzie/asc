'''
Redis service
'''''
import boto3
from .common import print_as_table

def add_subparsers(subparsers):
    redis_parser = subparsers.add_parser('redis', help='Redis service', description='Redis service',
                                   epilog='''Example: asc redis ls''')
    redis_parser.set_defaults(func=lambda args: redis_parser.print_help())
    redis_subparsers = redis_parser.add_subparsers(help='Description:', dest='subcommand')
    redis_list_parser = redis_subparsers.add_parser('ls')
    redis_list_parser.set_defaults(func=list_redis_instances)
    redis_tag_parser = redis_subparsers.add_parser('tag')
    redis_tag_parser.add_argument('name', help='Redis instance name')
    redis_tag_parser.add_argument('environment', help='Environment')
    redis_tag_parser.set_defaults(func=tag_redis)


def list_redis_instances(args):
    """
    List Redis instances
    """
    instance_list = []
    elasticache = boto3.client('elasticache')
    response = elasticache.describe_cache_clusters(ShowCacheNodeInfo=True)

    for cluster in response["CacheClusters"]:
        resource_tags = elasticache.list_tags_for_resource(
            ResourceName=cluster["ARN"])

        # Store the stack name
        for tag in resource_tags["TagList"]:
            if tag["Key"] == env_tag_key:
                env_tag_key = tag["Value"]
        instance = {}
        instance["Name"] = cluster["CacheClusterId"]
        instance["Endpoint"] = cluster["CacheNodes"][0]["Endpoint"]["Address"]
        instance["Type"] = cluster["CacheNodeType"]
        instance["Environment"] = env_tag_key
        instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i['Name'])
    print_as_table(instances)


def tag_redis(args):
    """
    Tag Redis instance
    """
    elasticache = boto3.client('elasticache')
