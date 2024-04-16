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

import logging
import operator
from collections import OrderedDict
from typing import Optional, List, Dict, Any
import sys
from tabulate import tabulate
from botocore.exceptions import UnauthorizedSSOTokenError
from boto3 import Session

SUBPARSER_REGISTRY = {}


def subparser_register(name: str):
    """
    Decorator for registering subparser functions.

    Args:
        name: The name of the command that the subparser will handle.
    """

    def decorator(func):
        SUBPARSER_REGISTRY[name] = func
        return func

    return decorator


def configure_logger(verbose_level: int):
    """
    Set up logging for the application

    Args:
        verbose_level: The level of verbosity to set for the logger.
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

    # Set up the logger
    if verbose_level > 0:
        logging.basicConfig(
            level=log_level,
            format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        )
    else:
        logging.basicConfig(level=log_level, format="%(message)s")

    logger = logging.getLogger("asc")
    logger.setLevel(log_level)

    # Set up the boto3 logger
    boto_logger = logging.getLogger("botocore")
    boto_logger.setLevel(boto_log_level)

    logger.info(
        "Log level set to %s, boto log level set to %s",
        log_level,
        boto_log_level,
    )


def create_boto_session(
    profile: Optional[str] = None, region: Optional[str] = None
) -> Session:
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
        session.client("sts").get_caller_identity()
        return session
    except UnauthorizedSSOTokenError as sso_err:
        logging.error("SSO Token Load Error: %s", sso_err)
        sys.exit(1)
    except Exception as e:
        # For all other exceptions, print the stack trace
        logging.error("Error: %s", e)
        sys.exit(1)


def print_as_table(
    items: List[Dict[str, Any]],
    key_order: Optional[List[str]] = None,
    sort_key: Optional[str] = None,
    sort_order: str = "asc",
):
    """
    Print a list of dicts as a table.

    Args:
        items: List of dictionaries containing the data to print.
        key_order: The order in which the keys should be arranged.
        sort_key: The key to sort the items by.

    Prints:
        A table representation of the provided data.
    """
    logging.debug("Printing %s items", len(items))
    if key_order:
        arranged_items = arrange_dict_keys(items, key_order)
    else:
        arranged_items = items

    if sort_key:
        try:
            arranged_items.sort(
                key=operator.itemgetter(sort_key), reverse=sort_order == "desc"
            )
        except KeyError:
            logging.error("Error: The key '%s' does not exist.", sort_key)
    table = tabulate(arranged_items, headers="keys")
    print(table)


def arrange_dict_keys(
    instance_list: List[Dict[str, Any]], key_order: Optional[List[str]]
) -> List[OrderedDict]:
    """
    Sorts all items in each dictionary in a list by the order of keys in
    key_order. If a key is not in any dictionary, it will not be included
    in the sorted list. If a key is in one dictionary but not in another,
    it is added with an empty value.

    Args:
        instance_list: The list of dictionaries to sort.
        key_order: The order in which the keys should be arranged.

    Returns:
        The sorted list of dictionaries.
    """
    logging.debug(
        "Arranging dictionary keys in the following order: %s", key_order
    )
    sorted_list = []

    # Iterate over each key in key_order
    for key in list(key_order):
        logging.debug("Checking for key: %s", key)
        # Check to see if the key is in any of the dictionaries
        if any(key in d for d in instance_list):
            # If the key exists in any of the dicts,
            # add it to any that are missing it
            for instance in instance_list:
                if key not in instance:
                    instance[key] = ""
        else:
            # Remove the key from the list if it doesn't exist
            key_order.remove(key)

    # Use key_order to sort the items in each dictionary
    key_order.reverse()
    for instance in instance_list:
        od = OrderedDict(instance)
        for key in key_order:
            od.move_to_end(key, last=False)

        sorted_list.append(od)

    return sorted_list


def apply_tags(
    instance: Dict[str, Any],
    instance_data: Dict[str, Any],
    displayed_tags_list: List[str],
) -> Dict[str, Any]:
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
    tags_key = "Tags" if "Tags" in instance_data else "TagList"

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
