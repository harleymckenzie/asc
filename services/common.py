import tabulate
import os
import configparser

def add_subparsers(subparsers):
    config_parser = subparsers.add_parser('configure', help='Configure asc', description='Configure asc',
                                          epilog='''Example: asc configure''')
    config_parser.set_defaults(func=configure)
    config_parser.add_argument('--environment-tag',
                               nargs='?',
                               help='AWS tag key used to identify environments')

def print_as_table(items):
    """
    Print a list of dicts as a table
    """
    print(tabulate.tabulate(items, headers="keys"))

def load_config():
    """
    Load configuration from ~/.asc/config if it exists
    and create it if it doesn't
    """
    config = configparser.RawConfigParser()
    config.add_section('asc')
    config_dir = os.path.expanduser('~/.asc')
    config_file = os.path.join(config_dir, 'config')

    # Read existing config if it exists and create an empty one if it doesn't
    if os.path.exists(config_file):
        config.read(config_file)
    else:
        generate_config(config, config_dir)

    # Return the config
    return config

def generate_config(config, config_dir):
    """
    Create an empty config file
    """
    if not os.path.exists(config_dir):
        os.makedirs(config_dir)
    with open(os.path.join(config_dir, 'config'), 'w') as configfile:
        config.write(configfile)


def configure(args):
    """
    Configure asc
    """
    config = configparser.RawConfigParser()
    config.add_section('asc')
    config_dir = os.path.expanduser('~/.asc')
    config_file = os.path.join(config_dir, 'config')

    env_tag_key = config.get('asc', 'env_tag_key', fallback='Environment')
    if args.environment_tag:
        env_tag_key = args.environment_tag
    else:
        env_tag_key = input(f'Environment tag [{env_tag_key}]: ') or env_tag_key
    config.set('asc', 'env_tag_key', env_tag_key)

    with open(os.path.join(config_dir, 'config'), 'w') as configfile:
        config.write(configfile)

    print('Configuration saved to ~/.asc/config')
