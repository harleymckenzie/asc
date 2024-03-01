"""
Testing utilities

This module contains utility functions for testing.
"""
from argparse import Namespace
import configparser


def setup_args(displayed_tags=None, profile=None, region=None):
    """
    Set up mock arguments for testing.
    """
    config = configparser.ConfigParser()
    config.add_section('asc')
    
    if displayed_tags:
        config.set('asc', 'displayed_tags', displayed_tags)

    args = Namespace(
        profile=profile,
        region=region,
        config=config,
    )

    return args
