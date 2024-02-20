"""
This module contains functions to interact with the Systems Manager service.
It provides functionality to list and connect to managed instances, and to
list and execute commands on managed instances.
"""
from ..common import subparser_register, print_as_table

@subparser_register('ssm')
def add_subparsers(subparsers, global_parser):
    """
    Adds subparsers for the SSM service.
    """
    ssm_parser = subparsers.add_parser(
        "ssm",
        help="Systems Manager service",
        description="Systems Manager service",
        parents=[global_parser],
    )
    ssm_parser.set_defaults(func=lambda args: ssm_parser.print_help())
    ssm_subparsers = ssm_parser.add_subparsers(
        help="Description:",
        dest="subcommand"
    )

    # Parameter Store Subparser
    add_parameter_store_parser(ssm_subparsers, global_parser)


def add_parameter_store_parser(ssm_subparsers, global_parser):
    """
    Adds a subparser for Parameter Store related commands.
    """
    ssm_parameter_parser = ssm_subparsers.add_parser(
        "parameter",
        help="Systems Manager Parameter Store subcommands",
        description="Systems Manager Parameter Store subcommands",
        epilog="""Example: asc ssm parameter ls""",
        parents=[global_parser],
    )
    ssm_parameter_parser.set_defaults(
        func=lambda args: ssm_parameter_parser.print_help()
    )
    ssm_parameter_subparsers = ssm_parameter_parser.add_subparsers(
        help="Description:",
        dest="subcommand"
    )

    # List Parameters Subcommand
    ssm_parameter_list_parser = ssm_parameter_subparsers.add_parser(
        "ls",
        help="List parameters",
        description="List parameters in the Parameter Store",
        epilog="""Example: asc ssm parameter ls""",
        parents=[global_parser],
    )
    ssm_parameter_list_parser.set_defaults(func=list_ssm_parameters)
    ssm_parameter_list_parser.add_argument(
        "--values", action="store_true", help="Include parameter values"
    )
    ssm_parameter_list_parser.add_argument(
        "--decrypt", action="store_true", help="Decrypt secure string parameters"
    )


def list_ssm_parameters(args):
    """
    List parameters in the Systems Manager Parameter Store.

    Args:
        args: Arguments passed by the user, including configuration details
              and options such as displaying endpoints.

    Prints:
        A table displaying the details of all parameters in the Parameter Store.
    """
    ssm_client = args.session.client("ssm")
    parameter_list = []

    try:
        response = ssm_client.describe_parameters()
    except Exception as e:
        print(f"Error listing parameters: {e}")
        exit(1)

    for instance in response["Parameters"]:
        instance = {
            "Name": instance["Name"],
            "Type": instance["Type"],
            "LastModifiedDate": instance["LastModifiedDate"]
        }
        
        parameter_list.append(instance)

    if args.values:
        parameter_list = get_ssm_parameters(args, parameter_list)

    print_as_table(parameter_list)


def get_ssm_parameters(args, parameter_list):
    """
    Get the values of the parameters in the Systems Manager Parameter Store.

    Args:
        args: Arguments passed by the user, including configuration details
              and options such as displaying endpoints.
        parameter_list: List of parameters to get values for.

    Returns:
        List of parameters with values.
    """
    ssm_client = args.session.client("ssm")

    try:
        response = ssm_client.get_parameters(
            Names=[parameter["Name"] for parameter in parameter_list],
            WithDecryption=args.decrypt
        )
    except Exception as e:
        print(f"Error getting parameter values: {e}")
        exit(1)

    # Update the parameter list with the parameter values
    # Set value to '*****' for SecureString parameters if --decrypt is not specified
    for parameter in parameter_list:
        for parameter_data in response["Parameters"]:
            if parameter["Name"] == parameter_data["Name"]:
                if parameter_data["Type"] == "SecureString" and not args.decrypt:
                    parameter["Value"] = "*****"
                else:
                    parameter["Value"] = parameter_data["Value"]
                break

    return parameter_list
