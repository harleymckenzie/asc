#!/usr/bin/env python
"""
'asc' is a simplified version of the AWS CLI.
"""
import argparse
import boto3
from asc.common import SUBPARSER_REGISTRY, load_config
from asc.services import asg, ec2, rds, redis


def main():
    """
    Main function
    """
    # Main parser
    # If no arguments are specified, print help
    parser = argparse.ArgumentParser(
        prog='asc',
        description='AWS Simple CLI (asc)',
        epilog='Example: asc ec2 ls'
    )
    parser.set_defaults(func=lambda args: parser.print_help())
    parser.add_argument(
        '--tags', '-t', help='Comma-separated tags to display in output.',
        type=str
    )
    parser.add_argument(
        '--profile', '-p', nargs='?', 
        help='AWS profile to use.',
        dest='profile'
    )
    parser.add_argument(
        '--region', nargs='?', help='AWS region to use.', dest='region'
    )

    subparsers = parser.add_subparsers(
        help='description', metavar='subcommand', dest='subcommand'
    )

    # Global parser
    # This parser will be used by all subparsers
    global_parser = argparse.ArgumentParser(add_help=False)
    group = global_parser.add_argument_group('global arguments')
    group.add_argument(
        '--profile', '-p', nargs='?',
        help='AWS profile to use.',
        dest='global_profile'
    )
    group.add_argument(
        '--region', nargs='?', help='AWS region to use.', dest='global_region'
    )

    for name, add_subparser_func in SUBPARSER_REGISTRY.items():
        add_subparser_func(subparsers, global_parser)

    args = parser.parse_args()

    # Load configuration
    args.config = load_config()

    # Combine tags from the config and command line
    if "displayed_tags" in args.config["asc"] and args.tags:
        args.config.set(
            'asc', 'displayed_tags',
            f"{args.config.get('asc', 'displayed_tags')},{args.tags}"
        )

    # Set up AWS session
    session_params = setup_session(args)

    try:
        args.session = boto3.Session(**session_params)
    except Exception as e:
        print(f"Failed to create AWS session: {e}")
        exit(1)

    args.func(args)


def setup_session(args):
    """
    Set up AWS session
    """
    session_params = {}

    if args.global_profile:
        session_params["profile_name"] = args.global_profile
    elif args.profile:
        session_params["profile_name"] = args.profile

    if args.global_region:
        session_params["region_name"] = args.global_region
    elif args.region:
        session_params["region_name"] = args.region

    return session_params


if __name__ == "__main__":
    main()
