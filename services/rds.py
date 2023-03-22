'''
RDS service
'''
import boto3
from .common import print_as_table

def add_subparsers(subparsers):
    rds_parser = subparsers.add_parser('rds', help='RDS service', description='RDS service',
                                   epilog='''Example: asc rds ls''')
    rds_parser.set_defaults(func=lambda args: rds_parser.print_help())
    rds_subparsers = rds_parser.add_subparsers(help='Description:', dest='subcommand')
    rds_list_parser = rds_subparsers.add_parser('ls', help='List RDS instances', description='List RDS instances',
                                                epilog='''Example: asc rds ls''')
    rds_list_parser.set_defaults(func=list_rds_instances)

def list_rds_instances(args):
    """
    List RDS instances
    """
    instance_list = []
    rds = boto3.client('rds')
    response = rds.describe_db_instances()
    if "aurora-mysql" in [db["Engine"] for db in response["DBInstances"]]:
        cluster_response = rds.describe_db_clusters()

    for db in response["DBInstances"]:
        if "TagList" in db:
            # Store the stack name
            for tag in db["TagList"]:
                if tag["Key"] == env_tag_key:
                    env_tag_key = tag["Value"]

        # Create a dict for each instance
        instance = {}
        instance["Name"] = db["DBInstanceIdentifier"]
        instance["Endpoint"] = db["Endpoint"]["Address"]
        instance["Type"] = db["DBInstanceClass"]
        instance["State"] = db["DBInstanceStatus"]
        instance["Environment"] = env_tag_key if "env_tag_key" in locals(
        ) else ""
        instance_list.append(instance)

        if db["Engine"] == "aurora-mysql":
            for cluster in cluster_response["DBClusters"]:
                for cluster_member in cluster["DBClusterMembers"]:
                    if cluster_member["DBInstanceIdentifier"] == db["DBInstanceIdentifier"]:
                        instance["Role"] = "Writer" if cluster_member["IsClusterWriter"] else "Reader"
                        instance["Endpoint"] = cluster["Endpoint"] if cluster_member["IsClusterWriter"] else cluster["ReaderEndpoint"]

    instances = sorted(instance_list, key=lambda i: i['Name'])
    print_as_table(instances)
