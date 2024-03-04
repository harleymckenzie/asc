"""
This module contains functions for working with the configuration
of the application.

Functions:
- initialise: Initialise the configuration.
- load: Load the configuration.
- setup: Run initial configuration setup for the application.
"""
import os
import sys
import configparser
from core.common import subparser_register

config_path = os.path.expanduser("~/.asc/config")


@subparser_register("configuration")
def add_subparsers(subparsers, global_parser):
    """
    Add subparsers for common commands.

    Args:
        subparsers: The subparsers object from the main parser.
        global_parser: The global parser containing common arguments.
    """
    config_parser = subparsers.add_parser(
        "configure",
        help="Configure asc",
        description="Configure asc",
        epilog="""Example: asc configure""",
        parents=[global_parser],
    )
    config_parser.set_defaults(func=update_config)


def initialise_config(tags=None):
    """
    Checks to see if the configuration file exists:
    If it does not exist, the setup function is called.
    If it does exist, the load function is called.

    Args:
        tags: Additional tags to include in displayed tags.

    Returns:
        The configuration object.
    """
    config = configparser.ConfigParser()

    if not os.path.exists(config_path):
        setup_config(config, initial_setup=True)
    config = load_config(config, tags)

    return config


def load_config(config, tags=None):
    """
    Load the configuration file and return the configuration object.
    If tags are provided, they are added to the displayed tags.

    Args:
        config: The configuration object.
        tags: Additional tags to include in displayed tags.

    Returns:
        The configuration object.
    """
    config.read(config_path)
    if tags:
        displayed_tags = config.get("asc", "displayed_tags")
        config.set("asc", "displayed_tags", f"{displayed_tags},{tags}")

    return config


def update_config(args):
    """
    Re-runs the configuration setup, priving the args.config configparser
    object.

    Args:
        args: The argparse args object.

    Returns:
        None
    """
    config = args.config
    setup_config(config)


def setup_config(config, initial_setup=False):
    """
    Run initial configuration setup for the application.
    If called for the first time, provided input + 'Name' is set.
    If called via 'configure' command, sets input value or existing as default.

    Args:
        config: The configuration object.
        initial_setup: If True, the setup is for the initial configuration.

    Returns:
        The configuration object.
    """
    print("\nConfiguration Setup\n" "-------------------")

    default_tags = "Name" if initial_setup else config["asc"]["displayed_tags"]

    try:
        tags_input = input(
            "Provide a comma separated list of tags to display [" +
            f"{default_tags}]:"
        ).strip() or default_tags
    except KeyboardInterrupt:
        print("\nConfiguration setup cancelled.")
        sys.exit(0)

    # If tags_input doens't contain 'Name', display a warning
    if "Name" not in tags_input:
        print(
            "WARNING: 'Name' is not included in the displayed tags. "
            "This may make it difficult to identify resources."
        )
    config["asc"] = {"displayed_tags": tags_input}
    save_config(config)


def save_config(config):
    """
    Save the configuration object to the configuration file.

    Args:
        config: The configuration object.

    Returns:
        None
    """
    try:
        with open(config_path, "w") as configfile:
            config.write(configfile)
        print("Configuration saved.")
    except Exception as e:
        print(f"Error saving configuration: {e}")
