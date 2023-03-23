"""
RDS service
"""
import boto3
from .common import print_as_table, load_config


def add_subparsers(subparsers):
    """
    Add subparsers for RDS commands
    """
    rds_parser = subparsers.add_parser('rds', help='RDS service', description='RDS service',
                                       epilog='''Example: asc rds ls''')
    rds_parser.set_defaults(func=lambda args: rds_parser.print_help())
    rds_subparsers = rds_parser.add_subparsers(help='Description:', dest='subcommand')

    # RDS list subcommand
    rds_list_parser = rds_subparsers.add_parser('ls', help='List RDS instances', description='List RDS instances',
                                                epilog='''Example: asc rds ls''')
    rds_list_parser.set_defaults(func=list_rds_instances)


def list_rds_instances(args):
    """
    List RDS instances
    """
    instance_list = []
    config = load_config()
    env_tag_key = config.get('asc', 'env_tag_key', fallback='Environment')
    rds = boto3.client('rds')
    response = rds.describe_db_instances()

    # Only call describe_db_clusters if there are Aurora instances
    if "aurora-mysql" in [db["Engine"] for db in response["DBInstances"]]:
        cluster_response = rds.describe_db_clusters()

    for db in response["DBInstances"]:
        instance = {"Name": db["DBInstanceIdentifier"], "Endpoint": db["Endpoint"]["Address"],
                    "Type": db["DBInstanceClass"], "State": db["DBInstanceStatus"]}

        # Store the environment tag in the instance dict if it exists
        if "TagList" in db:
            # Store the stack name
            for tag in db["TagList"]:
                if tag["Key"] == env_tag_key:
                    instance["Environment"] = tag["Value"]

        # Confirm whether the DB instance is a reader or writer
        if "aurora-mysql" in db["Engine"]:
            for cluster in cluster_response["DBClusters"]:
                for member in cluster["DBClusterMembers"]:
                    if member["DBInstanceIdentifier"] == db["DBInstanceIdentifier"]:
                        instance["Role"] = "Writer" if member["IsClusterWriter"] else "Reader"
                        instance["Endpoint"] = cluster["Endpoint"] if member["IsClusterWriter"] else cluster["ReaderEndpoint"]

        instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i['Name'])
    print_as_table(instances)
