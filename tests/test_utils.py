"""
Testing utilities

This module contains utility functions for testing.
"""
import configparser
from configparser import ConfigParser
from argparse import ArgumentParser, Namespace
from core.clidriver import setup_global_parser, process_args


def setup_args(profile=None, region=None, displayed_tags=None, sort_by=None):
    """
    Set up mock arguments for testing.
    """
    config = configparser.ConfigParser()
    config.add_section('asc')
    args = Namespace(
        profile=profile,
        region=region,
        config=config,
        sort_by=sort_by
    )
    
    if displayed_tags:
        config.set('asc', 'displayed_tags', displayed_tags)

    return args


def setup_parser(add_subparser_func, arg_list):
    """
    Parses a list of argument strings for EC2 commands.

    Args:
        add_subparser_func: The function to add subparsers to the parser.
        arg_list: A list of argument strings, e.g., ['ec2', 'ls', '--sort-by', 'Name']

    Returns:
        The parsed arguments namespace.
    """
    parser = ArgumentParser()
    subparsers = parser.add_subparsers()
    global_parser = setup_global_parser()  # Setup global parser options
    add_subparser_func(subparsers, global_parser)  # Add EC2 subparsers
    
    args = parser.parse_args(arg_list)
    args = process_args(args)
    return args


def setup_config(displayed_tags=None):
    """
    Set up a configparser.ConfigParser instance with the given displayed tags.
    """
    config = ConfigParser()
    config.add_section('asc')
    if displayed_tags:
        config.set('asc', 'displayed_tags', displayed_tags)
    return config
