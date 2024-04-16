"""
This module contains functions to interact with the Systems Manager service.
It provides functionality to list and connect to managed instances, and to
list and execute commands on managed instances.
"""

import json
import signal
import contextlib
import errno
import subprocess
import logging
from typing import Any
from ..common import subparser_register, create_boto_session, print_as_table


@subparser_register("ssm")
def add_subparsers(subparsers: Any, global_parser: Any) -> None:
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
        help="", metavar="subcommand", dest="subcommand"
    )

    # Parameter Store Subparser
    add_parameter_store_parser(ssm_subparsers, global_parser)
    add_session_manager_parser(ssm_subparsers, global_parser)


def add_parameter_store_parser(
    ssm_subparsers: Any, global_parser: Any
) -> None:
    """
    Adds a subparser for Parameter Store related commands.
    """
    ssm_parameter_parser = ssm_subparsers.add_parser(
        "parameter",
        aliases=["param"],
        help="Systems Manager Parameter Store subcommands",
        description="Systems Manager Parameter Store subcommands",
        epilog="""Example: asc ssm parameter ls""",
        parents=[global_parser],
    )
    ssm_parameter_parser.set_defaults(
        func=lambda args: ssm_parameter_parser.print_help()
    )
    ssm_parameter_subparsers = ssm_parameter_parser.add_subparsers(
        help="", metavar="subcommand", dest="subcommand"
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
        "--decrypt",
        action="store_true",
        help="Decrypt secure string parameters",
    )

    ssm_parameter_add_parser = ssm_parameter_subparsers.add_parser(
        "add",
        help="Add or update a parameter",
        description="Add a parameter to the Parameter Store",
        epilog="""Example: asc ssm parameter add""",
        parents=[global_parser],
    )
    ssm_parameter_add_parser.set_defaults(func=add_or_update_parameter)
    ssm_parameter_add_parser.add_argument(
        "name", help="Name of the parameter to add"
    )
    ssm_parameter_add_parser.add_argument(
        "value", help="Value of the parameter to add"
    )
    ssm_parameter_add_parser.add_argument(
        "--type",
        default="String",
        help="Type of the parameter to add",
        choices=["String", "SecureString"],
    )
    ssm_parameter_add_parser.add_argument(
        "--overwrite", action="store_true", help="Overwrite existing parameter"
    )

    ssm_parameter_remove_parser = ssm_parameter_subparsers.add_parser(
        "rm",
        help="Remove a parameter",
        description="Remove a parameter from the Parameter Store",
        epilog="""Example: asc ssm parameter remove""",
        parents=[global_parser],
    )
    ssm_parameter_remove_parser.set_defaults(func=remove_ssm_parameter)

    ssm_parameter_remove_parser.add_argument(
        "name", help="Name of the parameter to remove"
    )


def add_session_manager_parser(
    ssm_subparsers: Any, global_parser: Any
) -> None:
    """
    Adds a subparser for Session Manager related commands.
    """
    ssm_session_parser = ssm_subparsers.add_parser(
        "session",
        aliases=["connect"],
        help="Connect to a managed instance",
        description="Connect to a managed instance using Session Manager",
        epilog="""Example: asc ssm session i-1234567890abcdef0""",
        parents=[global_parser],
    )
    ssm_session_parser.set_defaults(func=connect_to_instance)
    ssm_session_parser.add_argument(
        "instance_id", help="ID of the instance to connect to"
    )


def list_ssm_parameters(args: Any) -> None:
    """
    List parameters in the Systems Manager Parameter Store.

    Args:
        args: Arguments passed by the user, including configuration details
              and options such as displaying endpoints.

    Prints:
        A table displaying the details of all parameters in the Parameter Store.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    ssm_client = session.client("ssm")
    parameter_list = []

    try:
        response = ssm_client.describe_parameters(MaxResults=50)
    except Exception as e:
        logging.error("Error listing parameters: %s", e)
        exit(1)

    for instance in response["Parameters"]:
        instance = {
            "Name": instance["Name"],
            "Type": instance["Type"],
            "Last Modified": instance["LastModifiedDate"],
        }

        parameter_list.append(instance)

    if args.values:
        parameter_list = get_ssm_parameters(args, parameter_list)

    print_as_table(parameter_list)


def get_ssm_parameters(args: Any, parameter_list: Any) -> Any:
    """
    Get the values of the parameters in the Systems Manager Parameter Store.

    Args:
        args: Arguments passed by the user, including configuration details
              and options such as displaying endpoints.
        parameter_list: List of parameters to get values for.

    Returns:
        List of parameters with values.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    ssm_client = session.client("ssm")

    # Call the get_parameters API with 10 parameters at a time
    for i in range(0, len(parameter_list), 10):
        chunk = parameter_list[i : i + 10]
        names = [param["Name"] for param in chunk]
        try:
            response = ssm_client.get_parameters(
                Names=names, WithDecryption=True
            )
            for param in response["Parameters"]:
                for chunk_param in chunk:
                    if chunk_param["Name"] == param["Name"]:
                        chunk_param["Value"] = param["Value"]
        except Exception as e:
            logging.error("Error getting parameter values: %s", e)
            exit(1)
    return parameter_list


def add_or_update_parameter(args: Any) -> None:
    """
    Add or update a parameter in the Systems Manager Parameter Store.

    Args:
        args: Arguments passed by the user, including parameter details.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    ssm_client = session.client("ssm")

    try:
        ssm_client.put_parameter(
            Name=args.name,
            Value=args.value,
            Type=args.type,
            Overwrite=args.overwrite,
        )
        print(
            f"Parameter {'updated' if args.overwrite else 'added'} successfully."
        )
    except ssm_client.exceptions.ParameterAlreadyExists:
        print(
            f"Parameter '{args.name}' already exists. Use --overwrite to update."
        )
        exit(1)
    except Exception as e:
        logging.error("Error adding/updating parameter: %s", e)
        exit(1)


def remove_ssm_parameter(args: Any) -> None:
    """
    Remove a parameter from the Systems Manager Parameter Store.

    Args:
        args: Arguments passed by the user, including the name of the parameter to remove.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    ssm_client = session.client("ssm")

    try:
        ssm_client.delete_parameter(Name=args.name)
        print(f"Parameter '{args.name}' removed successfully.")
    except Exception as e:
        logging.error("Error removing parameter: %s", e)
        exit(1)


@contextlib.contextmanager
def ignore_user_entered_signals():
    """
    Ignore signals sent by the user, such as Ctrl+C and Ctrl+Z.
    """
    signal_list = [signal.SIGINT, signal.SIGQUIT, signal.SIGTSTP]
    actual_signals = []
    for user_signal in signal_list:
        actual_signals.append(signal.signal(user_signal, signal.SIG_IGN))
    try:
        yield
    finally:
        for sig, user_signal in enumerate(signal_list):
            signal.signal(user_signal, actual_signals[sig])


def connect_to_instance(args: Any) -> None:
    """
    Connect to a managed instance using Session Manager.

    Args:
        args: Arguments passed by the user, including the instance ID and
              whether to connect using SSH.

    Prints:
        The command to connect to the instance using Session Manager or SSH.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    ssm_client = session.client("ssm")
    instance_id = args.instance_id
    response = ssm_client.start_session(Target=instance_id)
    session_id = response["SessionId"]
    region_name = session.region_name
    profile_name = session.profile_name if session.profile_name else ""
    endpoint_url = ssm_client.meta.endpoint_url

    try:
        with ignore_user_entered_signals():
            subprocess.check_call(
                [
                    "session-manager-plugin",
                    json.dumps(response),
                    region_name,
                    "StartSession",
                    profile_name,
                    json.dumps(dict(Target=instance_id)),
                    endpoint_url,
                ]
            )
        return 0
    except OSError as ex:
        if ex.errno == errno.ENOENT:
            ssm_client.terminate_session(SessionId=session_id)
            logging.error(
                "The session-manager-plugin executable could not be found."
            )


def is_plugin_installed() -> bool:
    """Check if the session-manager-plugin is installed."""
    try:
        subprocess.check_call(
            ["session-manager-plugin"],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False
