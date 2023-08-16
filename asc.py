#!/usr/bin/env python
"""
'asc' is a simplified version of the AWS CLI.
"""
import argparse
import boto3
from services import ec2, rds, asg, redis, common


def main():
    """
    Main function
    """
    # Main parser
    # If no arguments are specified, print help
    parser = argparse.ArgumentParser(prog='asc', description='AWS Simple CLI (asc)',
                                     epilog='''Example: asc ec2 ls''')
    parser.set_defaults(func=lambda args: parser.print_help())
    parser.add_argument('--tags', '-t', help='Comma-separated tags to display in output.',
                        type=str)
    parser.add_argument('--profile', '-p', nargs='?', help='AWS profile to use.', dest='profile')
    parser.add_argument('--region', nargs='?', help='AWS region to use.', dest='region')
    
    subparsers = parser.add_subparsers(help='description', metavar='subcommand', dest='subcommand')

    # Global parser
    # This parser will be used by all subparsers
    global_parser = argparse.ArgumentParser(add_help=False)
    group = global_parser.add_argument_group('global arguments')
    group.add_argument('--profile', '-p', nargs='?', help='AWS profile to use.', dest='global_profile')
    group.add_argument('--region', nargs='?', help='AWS region to use.', dest='global_region')
    
    for service in [common, ec2, rds, asg, redis]:
        service.add_subparsers(subparsers, global_parser)

    args = parser.parse_args()
    print(args)

    # Load configuration
    args.config = common.load_config()

    # Combine tags from the config and command line
    if "displayed_tags" in args.config["asc"] and args.tags:
        args.config.set('asc', 'displayed_tags', f"{args.config.get('asc', 'displayed_tags')},{args.tags}")

    # Set up AWS session
    # If a profile or region is specified in the main and global parsers, use the global parser's values
    session_params = {}
    if args.global_profile:
        session_params["profile_name"] = args.global_profile
    elif args.profile:
        session_params["profile_name"] = args.profile
    if args.global_region:
        session_params["region_name"] = args.global_region
    elif args.region:
        session_params["region_name"] = args.region

    try:
        args.session = boto3.Session(**session_params)
    except Exception as e:
        print(f"Failed to create AWS session: {e}")
        exit(1)

    args.func(args)

if __name__ == "__main__":
    main()
