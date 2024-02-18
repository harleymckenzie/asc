"""
Common functions for asc
"""
import os
import configparser
import tabulate


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
    config_parser.set_defaults(func=configure)
    config_parser.add_argument(
        "--add-tag",
        nargs="?",
        help="Add a tag to the list of defined tags that are displayed",
    )
    config_parser.add_argument(
        "--remove-tag",
        "--rm-tag",
        nargs="?",
        help="Remove a tag from the list of defined tags that are displayed",
    )


def print_as_table(items):
    """
    Print a list of dicts as a table.

    Args:
        items: List of dictionaries containing the data to print.

    Prints:
        A table representation of the provided data.
    """

    print(tabulate.tabulate(items, headers="keys"))


def load_config():
    """
    Load configuration from ~/.asc/config if it exists and create it if it doesn't.

    Returns:
        A configparser object with the loaded configuration.
    """
    config = configparser.RawConfigParser()
    config_dir = os.path.expanduser("~/.asc")
    config_file = os.path.join(config_dir, "config")

    # Ensure directory exists
    if not os.path.exists(config_dir):
        os.makedirs(config_dir)

    # Create a new config if it doesn't exist
    if os.path.exists(config_file):
        config.read(config_file)
    else:
        config.add_section("asc")
        config.set("asc", "displayed_tags", "Name")
        with open(config_file, "w") as configfile:
            config.write(configfile)

    if "displayed_tags" not in config["asc"]:
        config.set("asc", "displayed_tags")

    # Return the config
    return config


def configure(args):
    """
    Configure asc based on the given arguments.

    Args:
        args: The arguments received from the command-line input.

    Prints:
        A confirmation message indicating the configuration has been saved.
    """
    config = load_config()
    config_dir = os.path.expanduser("~/.asc")

    displayed_tags = config.get("asc", "displayed_tags", fallback="").split(",")

    if args.add_tag and args.add_tag not in displayed_tags:
        displayed_tags.append(args.add_tag)
    if args.remove_tag and args.remove_tag in displayed_tags:
        displayed_tags.remove(args.remove_tag)

    config.set("asc", "displayed_tags", ",".join(displayed_tags))

    with open(os.path.join(config_dir, "config"), "w") as configfile:
        config.write(configfile)

    print("Configuration saved to ~/.asc/config")
