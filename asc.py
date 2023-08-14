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
    global_parser = argparse.ArgumentParser(add_help=False)
    group = global_parser.add_argument_group('global arguments')
    group.add_argument('--profile', nargs='?', help='AWS profile to use.', )
    group.add_argument('--region', nargs='?', help='AWS region to use.')
    
    # Parse command line arguments
    # If no arguments are specified, print help
    parser = argparse.ArgumentParser(prog='asc', description='AWS Simple CLI (asc)',
                                     epilog='''Example: asc ec2 ls''', parents=[global_parser])
    parser.set_defaults(func=lambda args: parser.print_help())
    
    # Add optional arguments
    # parser.add_argument('--profile', nargs='?', help='AWS profile to use.')
    # parser.add_argument('--region', nargs='?', help='AWS region to use.')
    # Specify additional tags to display in output. This will append to the displayed_tags list
    parser.add_argument('--tags', '-t', help='Comma-separated tags to display in output.',
                        type=str)
    
    subparsers = parser.add_subparsers(help='description', metavar='subcommand')
    for service in [common, rds, asg, redis]:
        service.add_subparsers(subparsers)

    for service in [ec2]:
        service.add_subparsers(subparsers, global_parser)

    args = parser.parse_args()

    # Load configuration
    args.config = common.load_config()

    # If tags are specified in the config file as well as the command line, append the command line tags
    if "displayed_tags" in args.config["asc"] and args.tags:
        args.config.set('asc', 'displayed_tags', f"{args.config.get('asc', 'displayed_tags')},{args.tags}")

    # Set up AWS session
    session_params = {}
    if args.profile:
        session_params['profile_name'] = args.profile
    if args.region:
        session_params['region_name'] = args.region

    try:
        args.session = boto3.Session(**session_params)
    except Exception as e:
        print(f"Failed to create AWS session: {e}")
        exit(1)

    try:
        args.func(args)
    except Exception as e:
        # Again, refine this catch based on known exceptions your service functions might raise.
        print(f"Error executing function: {e}")

if __name__ == "__main__":
    main()
