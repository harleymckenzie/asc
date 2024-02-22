#!/usr/bin/env python
"""
'asc' is a simplified version of the AWS CLI.
"""
import argparse
from asc import common
from asc.services import asg, ec2, rds, redis, ssm


def arg_parser():
    """
    Create the main parser
    """
    parser = argparse.ArgumentParser(
        prog='asc',
        description='AWS Simple CLI (asc)',
        epilog='Example: asc ec2 ls',
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
    global_parser = setup_global_parser()

    return parser, subparsers, global_parser


def setup_global_parser():
    """
    Create the global parser
    """
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
    group.add_argument(
        '-v', '--verbose', action='count',
        help='Increase verbosity level. 0 = WARNING, 1 = INFO, 2 = DEBUG, 3 = DEBUG with boto3',
        dest='global_verbose'
    )

    return global_parser


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

    # Set up logging
    common.logger(args)

    # Load configuration
    args.config = common.init_config()

    # Combine tags from the config and command line
    if "displayed_tags" in args.config["asc"] and args.tags:
        args.config.set(
            'asc', 'displayed_tags',
            f"{args.config.get('asc', 'displayed_tags')},{args.tags}"
        )

    # Set up AWS session
    args.session = common.create_boto_session(args)

    args.func(args)


if __name__ == "__main__":
    main()
