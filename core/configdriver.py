"""
This module contains functions for working with the configuration
of the application.

Functions:
- initialise: Initialise the configuration.
- load: Load the configuration.
- setup: Run initial configuration setup for the application.
"""
import os
import configparser

config_path = os.path.expanduser("~/.asc/config")


def initialise(tags=None):
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
        setup(config, initial_setup=True)
    config = load(config, tags)

    return config


def load(config, tags=None):
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
        displayed_tags = config.get('asc', 'displayed_tags')
        config.set('asc', 'displayed_tags', f"{displayed_tags},{tags}")

    return config


def setup(config, initial_setup=False):
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
    print("Configuration Setup:")
    if initial_setup:
        tags_input = input("Enter any tags you would like to displayed in "
                           "resource outputs, as comma separated values "
                           "(automatically includes 'Name'): ").strip()
        tags = "Name" if not tags_input else f"Name,{tags_input}"
        config['asc'] = {'displayed_tags': tags}
    else:
        current_tags = config['asc']['displayed_tags']
        print(f"Current displayed tags: {current_tags}")
        tags_input = input("Enter new displayed tags as comma separated values "
                           "or press enter to keep existing: ").strip()
        if tags_input:
            config['asc']['displayed_tags'] = tags_input

    try:
        with open(config_path, 'w') as configfile:
            config.write(configfile)
        print("Configuration saved.")
    except Exception as e:
        print(f"Error saving configuration: {e}")
