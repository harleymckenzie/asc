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


@subparser_register('common')
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
    config_parser.set_defaults(func=setup_config)
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


def setup_config(config):
    """
    Run initial configuration setup for the application.
    If called for the first time, it sets default values.
    If called via 'configure' command, it allows updating existing values.
    """
    print("Configuration Setup:")

    current_tags = config.get('asc', 'displayed_tags', fallback="Name")
    new_tags = input(
        f"Enter displayed tags (current: {current_tags}, leave blank to keep): "
        ).strip()
    config.set('asc', 'displayed_tags', new_tags if new_tags else current_tags)

    return config


def init_config():
    """Load the configuration or initialize it if it doesn't exist."""
    config_dir = os.path.expanduser("~/.asc")
    config_file = os.path.join(config_dir, "config")
    config = configparser.ConfigParser()

    if not os.path.exists(config_dir):
        os.makedirs(config_dir)

    if not os.path.exists(config_file):
        config.add_section("asc")
        config = setup_config(config)
        with open(config_file, "w") as configfile:
            config.write(configfile)
        print("Initial configuration saved.")
    else:
        config.read(config_file)

    return config


def print_as_table(items):
    """
    Print a list of dicts as a table.

    Args:
        items: List of dictionaries containing the data to print.

    Prints:
        A table representation of the provided data.
    """

    print(tabulate.tabulate(items, headers="keys"))


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
