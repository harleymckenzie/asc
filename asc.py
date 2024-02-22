#!/usr/bin/env python
"""
'asc' is a simplified version of the AWS CLI.
"""
import logging
import argparse
import boto3
from asc import common
from asc.services import asg, ec2, rds, redis, ssm


def arg_parser():
    """
    Create the main parser
    """
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
    parser.add_argument('-v', '--verbose', action='count', default=0,
                        help='Increase verbosity level. Use -v for INFO level and -vv for DEBUG level.')

    subparsers = parser.add_subparsers(
        help='', metavar='subcommand', dest='subcommand'
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

    return parser, subparsers, global_parser


def main():
    """
    Main function
    """
    # Main parser
    # If no arguments are specified, print help
    parser, subparsers, global_parser = arg_parser()

    for _, add_subparser_func in common.SUBPARSER_REGISTRY.items():
        add_subparser_func(subparsers, global_parser)

    args = parser.parse_args()

    # Load configuration
    args.config = common.init_config()

    # Combine tags from the config and command line
    if "displayed_tags" in args.config["asc"] and args.tags:
        args.config.set(
            'asc', 'displayed_tags',
            f"{args.config.get('asc', 'displayed_tags')},{args.tags}"
        )

    # Set up AWS session
    args.session = setup_session(args)

    args.func(args)


def setup_session(args):
    """
    Set up AWS session
    """
    session_params = {}

    # Handle profile and region arguments on both the regular and global parsers
    profile = getattr(args, 'global_profile', None) or getattr(args, 'profile', None)
    region = getattr(args, 'global_region', None) or getattr(args, 'region', None)

    if profile:
        session_params["profile_name"] = profile

    if region:
        session_params["region_name"] = region

    try:
        args.session = boto3.Session(**session_params)
        logging.debug("Using AWS profile: %s", args.session.profile_name)
    except Exception as e:
        print(f"Failed to create AWS session: {e}")
        exit(1)
    return args.session


if __name__ == "__main__":
    logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s')
    main()
