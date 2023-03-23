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
    # Parse command line arguments
    # If no arguments are specified, print help
    parser = argparse.ArgumentParser(prog='asc', description='AWS Simple CLI (asc)',
                                     epilog='''Example: asc ec2 ls''')
    parser.set_defaults(func=lambda args: parser.print_help())
    # Add optional arguments
    parser.add_argument('--profile', nargs='?', help='AWS profile to use')
    parser.add_argument('--region', nargs='?', help='AWS region to use')

    subparsers = parser.add_subparsers(help='description', metavar='subcommand')

    for service in [common, ec2, rds, asg, redis]:
        service.add_subparsers(subparsers)

    args = parser.parse_args()

    # If a profile and/or region is specified, use it
    if args.profile:
        boto3.setup_default_session(profile_name=args.profile)
    if args.region:
        boto3.setup_default_session(region_name=args.region)

    # Run the function for the specified service
    args.func(args)


if __name__ == "__main__":
    main()
