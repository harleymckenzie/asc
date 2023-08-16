"""
RDS service
"""
import boto3
from .common import print_as_table


def add_subparsers(subparsers, global_parser):
    """
    Add subparsers for RDS commands
    """
    rds_parser = subparsers.add_parser('rds', help='RDS service', description='RDS service',
                                       epilog='''Example: asc rds ls''', parents=[global_parser])
    rds_parser.set_defaults(func=lambda args: rds_parser.print_help())
    rds_subparsers = rds_parser.add_subparsers(help='Description:', dest='subcommand')

    # RDS list subcommand
    rds_list_parser = rds_subparsers.add_parser('ls', help='List RDS instances', description='List RDS instances',
                                                epilog='''Example: asc rds ls''', parents=[global_parser])
    rds_list_parser.add_argument('--endpoint', '-e', help='Display endpoint in output.', action='store_true')
    rds_list_parser.set_defaults(func=list_rds_instances)


def list_rds_instances(args):
    """
    List RDS instances
    """
    instance_list = []
    rds_client = args.session.client('rds')
    response = rds_client.describe_db_instances()

    # Store tags to display in the output if they've been set in the config
    if "displayed_tags" in args.config["asc"]:
        displayed_tags_list = args.config.get('asc', 'displayed_tags').split(',')
    # Set an empty list if the config hasn't been set
    else:
        displayed_tags_list = []

    # Only call describe_db_clusters if there are Aurora instances
    if "aurora-mysql" in [db["Engine"] for db in response["DBInstances"]]:
        cluster_response = rds_client.describe_db_clusters()

    for db in response["DBInstances"]:
        instance = {"Name": db["DBInstanceIdentifier"],
                    "Type": db["DBInstanceClass"], 
                    "State": db["DBInstanceStatus"]}

        # Add Endpoint if args.endpoint is set
        if args.endpoint:
            instance["Endpoint"] = db["Endpoint"]["Address"]

        # Add tags to instance dict
        for tag in db.get("TagList", []):
            if tag["Key"] in displayed_tags_list:
                instance[tag["Key"]] = tag["Value"]

        # Confirm whether the DB instance is a reader or writer to display in output
        if "aurora-mysql" in db["Engine"]:
            for cluster in cluster_response["DBClusters"]:
                for member in cluster["DBClusterMembers"]:
                    if member["DBInstanceIdentifier"] == db["DBInstanceIdentifier"]:
                        instance["Role"] = "Writer" if member["IsClusterWriter"] else "Reader"
                        # Add endpoint to instance dict only if --endpoint flag is set
                        if args.endpoint:
                            instance["Endpoint"] = cluster["Endpoint"] if member["IsClusterWriter"] else cluster["ReaderEndpoint"]

        instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i['Name'])
    print_as_table(instances)
