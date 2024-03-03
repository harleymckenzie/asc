"""
Common functions for asc

This module contains common functions and decorators used across the
application.

Functions:
- subparser_register: Decorator for registering subparser functions.
- add_subparsers: Add subparsers for common commands.
- init_config: Load the configuration or initialize it if it doesn't exist.
- print_as_table: Print a list of dicts as a table.
- apply_tags: Apply tags to an instance.
"""
import os
import logging
from tabulate import tabulate
from botocore.exceptions import UnauthorizedSSOTokenError
from boto3 import Session

SUBPARSER_REGISTRY = {}


def subparser_register(name):
    """
    Decorator for registering subparser functions.

    Args:
        name: The name of the command that the subparser will handle.
    """

    def decorator(func):
        SUBPARSER_REGISTRY[name] = func
        return func

    return decorator


def logger(verbose_level):
    """
    Set up logging for the application
    """

    if verbose_level == 1:
        log_level = "INFO"
        boto_log_level = "WARNING"
    elif verbose_level == 2:
        log_level = "DEBUG"
        boto_log_level = "INFO"
    elif verbose_level > 2:
        log_level = "DEBUG"
        boto_log_level = "DEBUG"
    else:
        log_level = "WARNING"
        boto_log_level = "ERROR"

    logging.basicConfig(
        level=log_level,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
    )
    logging.getLogger("botocore").setLevel(boto_log_level)


def create_boto_session(profile=None, region=None):
    """
    Set up AWS session.

    Args:
        profile: The profile to use for the session.
        region: The region to use for the session.

    Returns:
        The AWS session.
    """
    session_params = {}

    if profile:
        session_params["profile_name"] = profile
    if region:
        session_params["region_name"] = region

    try:
        session = Session(**session_params)
        session.client('sts').get_caller_identity()
        return session
    except UnauthorizedSSOTokenError as sso_err:
        print(f"SSO Token Load Error: {sso_err}")
        exit(1)
    except Exception as e:
        # For all other exceptions, print the stack trace
        print(f"Error: {e}")
        exit(1)


def print_as_table(items):
    """
    Print a list of dicts as a table.

    Args:
        items: List of dictionaries containing the data to print.

    Prints:
        A table representation of the provided data.
    """
    table = tabulate(items, headers="keys")
    print(table)


def apply_tags(instance, instance_data, displayed_tags_list):
    """
    Check to see if items in displayed tags list are in the instance data
    The "Name" tag should automatically be displayed

    Args:
        instance: The instance to apply tags to.
        instance_data: The data for the instance.
        displayed_tags_list: The list of tags to display.

    Returns:
        The instance with the tags applied.
    """
    # Determine whether to use 'Tags' or 'TagList'
    tags_key = 'Tags' if 'Tags' in instance_data else 'TagList'

    for tag_dict in instance_data.get(tags_key, []):
        tag_key = tag_dict.get("Key")
        tag_value = tag_dict.get("Value")

        if tag_key in displayed_tags_list:
            # Special handling for "Name" tag to prepend it
            if tag_key == "Name":
                instance = {tag_key: tag_value, **instance}
            else:
                instance[tag_key] = tag_value

    return instance
