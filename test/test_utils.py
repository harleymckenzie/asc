"""
This module provides utilities for testing the asc.py script.

Functions:
- setup_args(displayed_tags=None): Set up mock arguments for testing.
- run_and_capture_output(test_func, args): Run a test function and capture its output.
- check_output_for_tags(output, displayed_tags, unspecified_tags):
    Check if the output contains the specified tags and not the unspecified ones.
- mock_ec2_client(): Provide a mocked EC2 client using Moto.
- mock_rds_client(): Provide a mocked RDS client using Moto.
"""

from unittest.mock import MagicMock, patch
import configparser
import boto3
from moto import mock_ec2, mock_rds, mock_ssm
import pytest

def setup_args(displayed_tags=None):
    """
    Set up mock arguments for testing.
    """
    args = MagicMock()
    args.session = boto3.Session(region_name="eu-west-1")
    args.config = configparser.RawConfigParser()
    args.config.add_section("asc")
    if displayed_tags is not None:
        args.config.set("asc", "displayed_tags", displayed_tags)
    return args

def run_and_capture_output(test_func, args):
    """
    Run a test function and capture its output.
    """
    with patch('builtins.print') as mock_print:
        test_func(args)
        return [call_arg[0][0] for call_arg in mock_print.call_args_list]

def check_output_for_tags(output, displayed_tags, unspecified_tags):
    """
    Check if the output contains the specified tags and not the unspecified ones.
    """
    for tag in displayed_tags.split(",") if displayed_tags else []:
        assert tag in output
    for tag in unspecified_tags:
        assert tag not in output

@pytest.fixture
def mock_ec2_client():
    """
    Provide a mocked EC2 client using Moto.
    """
    with mock_ec2():
        yield boto3.client("ec2", region_name="eu-west-1")

@pytest.fixture
def mock_rds_client():
    """
    Provide a mocked RDS client using Moto.
    """
    with mock_rds():
        yield boto3.client("rds", region_name="eu-west-1")

@pytest.fixture
def mock_ssm_client():
    """
    Provide a mocked SSM client using Moto.
    """
    with mock_ssm():
        yield boto3.client("ssm", region_name="eu-west-1")