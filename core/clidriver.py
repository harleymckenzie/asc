from argparse import ArgumentParser
from importlib.metadata import version
from core import common, configdriver
from core.services import asg, ec2, rds, redis, ssm


def arg_parser():
    """
    Create the main parser
    """
    parser = ArgumentParser(
        prog='asc',
        description='AWS Simple CLI (asc)',
        epilog='Example: asc ec2 ls'
    )
    parser.set_defaults(func=lambda args: parser.print_help())
    parser.add_argument(
        '--version', action='version', version=version('asc')
    )
    group = parser.add_argument_group('global arguments')
    group.add_argument(
        '--tags', '-t', help='Comma-separated tags to display in output.',
        type=str
    )
    group.add_argument(
        '--profile', '-p', nargs='?',
        help='AWS profile to use.',
        dest='profile'
    )
    group.add_argument(
        '--region', nargs='?', help='AWS region to use.', dest='region'
    )
    group.add_argument(
        '-v', action='count', default=0,
        help='Increase verbosity of output.', dest='verbose'
    )

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
    global_parser = ArgumentParser(add_help=False)
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
        '-v', action='count',
        help='Increase verbosity of output.',
        dest='global_verbose'
    )

    return global_parser


def process_args(args):
    """
    Process the arguments and global arguments
    """
    if hasattr(args, 'global_profile') and args.global_profile:
        args.profile = args.global_profile
    if hasattr(args, 'global_region') and args.global_region:
        args.region = args.global_region
    if hasattr(args, 'global_verbose') and args.global_verbose:
        args.verbose = args.global_verbose

    return args


def main():
    """
    Main function
    """
    # If no arguments are specified, print help
    parser, subparsers, global_parser = arg_parser()

    for _, add_subparser_func in common.SUBPARSER_REGISTRY.items():
        add_subparser_func(subparsers, global_parser)

    args = parser.parse_args()
    args = process_args(args)

    # Set up logging
    common.logger(args.verbose)

    # Load configuration
    args.config = configdriver.initialise(args.tags)

    args.func(args)


if __name__ == "__main__":
    main()
