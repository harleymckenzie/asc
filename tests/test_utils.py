"""
Testing utilities

This module contains utility functions for testing.
"""
import argparse
import configparser
import pytest
from unittest.mock import patch


def setup_args(displayed_tags=None, profile=None, region=None):
    """
    Set up mock arguments for testing.
    """
    config = configparser.ConfigParser()
    config.add_section('asc')
    
    if displayed_tags:
        config.set('asc', 'displayed_tags', displayed_tags)

    args = argparse.Namespace(
        profile=profile,
        region=region,
        config=config,
    )

    return args
