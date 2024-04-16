"""
Tests for the ASG service

This module contains tests for the ASG service.
"""

from unittest.mock import patch
import logging
import datetime
from dateutil.tz import tzutc
import pytest
from core.services.asg import (
    add_subparsers,
    list_autoscaling_groups,
    list_autoscaling_schedules,
    add_autoscaling_schedule,
)
from .test_utils import setup_parser, setup_config


@pytest.mark.parametrize(
    "arg_list, displayed_tags, expected_output",
    [
        (
            ["asg", "ls"],
            "Environment",
            (
                "ASG Name              Min    Max    Desired  Environment\n"
                "------------------  -----  -----  ---------  -------------\n"
                "asg-web-production      2      8          2  production\n"
                "asg-web-staging         1      1          1\n"
            ),
        ),  # Displayed tags
        (
            ["asg", "ls", "--sort-by", "Min"],
            None,  # Displayed tags
            (
                "ASG Name              Min    Max    Desired\n"
                "------------------  -----  -----  ---------\n"
                "asg-web-staging         1      1          1\n"
                "asg-web-production      2      8          2\n"
            ),
        ),
    ],
    ids=["List ASGs", "List ASGs sorted by StartTime"],
)
@patch("core.services.asg.create_boto_session")
def test_list_autoscaling_groups(
    mock_create_boto_session, arg_list, displayed_tags, expected_output, capsys
):
    args = setup_parser(add_subparsers, arg_list)
    args.config = setup_config(displayed_tags)

    mock_session = mock_create_boto_session.return_value
    mock_client = mock_session.client.return_value
    mock_client.describe_auto_scaling_groups.return_value = {
        "AutoScalingGroups": [
            {
                "AutoScalingGroupName": "asg-web-production",
                "DesiredCapacity": 2,
                "MaxSize": 8,
                "MinSize": 2,
                "Tags": [{"Key": "Environment", "Value": "production"}],
            },
            {
                "AutoScalingGroupName": "asg-web-staging",
                "DesiredCapacity": 1,
                "MaxSize": 1,
                "MinSize": 1,
                "Tags": [{}],
            },
        ]
    }
    list_autoscaling_groups(args)
    out, _ = capsys.readouterr()
    print("\n" + out)
    assert out.strip() == expected_output.strip()

    # Verify that external calls were made as expected
    mock_create_boto_session.assert_called_once_with(
        profile=args.profile, region=args.region
    )
    mock_client.describe_auto_scaling_groups.assert_called_once()


@pytest.mark.parametrize(
    "arg_list, displayed_tags, expected_output",
    [
        (
            ["asg", "schedule", "ls"],
            "Environment",  # Displayed tags
            (
                "ASG Name            Name                   Start Time                   Min  Max\n"
                "------------------  ---------------------  -------------------------  -----  -----\n"
                "asg-web-production  Production schedule 1  2024-04-10 10:00:00+00:00      2\n"
                "asg-web-production  Production schedule 2  2024-04-03 00:00:00+00:00      3  8\n"
                "asg-web-staging     Staging schedule 1     2024-04-03 10:00:00+00:00      1\n"
            ),
        ),
        (
            ["asg", "schedule", "ls", "--sort-by", "Min"],
            None,  # Displayed tags
            (
                "ASG Name            Name                   Start Time                   Min  Max\n"
                "------------------  ---------------------  -------------------------  -----  -----\n"
                "asg-web-staging     Staging schedule 1     2024-04-03 10:00:00+00:00      1\n"
                "asg-web-production  Production schedule 1  2024-04-10 10:00:00+00:00      2\n"
                "asg-web-production  Production schedule 2  2024-04-03 00:00:00+00:00      3  8\n"
            ),
        ),
        (
            ["asg", "schedule", "ls", "--sort-by", "Start Time"],
            None,  # Displayed tags
            (
                "ASG Name            Name                   Start Time                   Min  Max\n"
                "------------------  ---------------------  -------------------------  -----  -----\n"
                "asg-web-production  Production schedule 2  2024-04-03 00:00:00+00:00      3  8\n"
                "asg-web-staging     Staging schedule 1     2024-04-03 10:00:00+00:00      1\n"
                "asg-web-production  Production schedule 1  2024-04-10 10:00:00+00:00      2\n"
            ),
        ),
    ],
    ids=["List ASGs",
         "List ASGs sorted by Min",
         "List ASGS sorted by StartTime"],
)
@patch("core.services.asg.create_boto_session")
def test_list_autoscaling_schedules(
    mock_create_boto_session, arg_list, displayed_tags, expected_output, capsys
):
    args = setup_parser(add_subparsers, arg_list)
    args.config = setup_config(displayed_tags)

    mock_session = mock_create_boto_session.return_value
    mock_client = mock_session.client.return_value
    mock_client.describe_scheduled_actions.return_value = {
        "ScheduledUpdateGroupActions": [
            {
                "AutoScalingGroupName": "asg-web-production",
                "MinSize": 2,
                "ScheduledActionName": "Production schedule 1",
                "StartTime": datetime.datetime(
                    2024, 4, 10, 10, 0, tzinfo=tzutc()
                ),
                "Time": datetime.datetime(2024, 4, 10, 10, 0, tzinfo=tzutc()),
            },
            {
                "AutoScalingGroupName": "asg-web-production",
                "MinSize": 3,
                "MaxSize": 8,
                "ScheduledActionName": "Production schedule 2",
                "StartTime": datetime.datetime(
                    2024, 4, 3, 0, 0, tzinfo=tzutc()
                ),
                "Time": datetime.datetime(2024, 4, 3, 0, 0, tzinfo=tzutc()),
            },
            {
                "AutoScalingGroupName": "asg-web-staging",
                "MinSize": 1,
                "ScheduledActionName": "Staging schedule 1",
                "StartTime": datetime.datetime(
                    2024, 4, 3, 10, 0, tzinfo=tzutc()
                ),
                "Time": datetime.datetime(2024, 4, 3, 10, 0, tzinfo=tzutc()),
            },
        ]
    }

    list_autoscaling_schedules(args)
    out, _ = capsys.readouterr()
    logging.debug("\n %s", out)
    print("\n" + out)
    assert out.strip() == expected_output.strip()

    # Verify that external calls were made as expected
    mock_create_boto_session.assert_called_once_with(
        profile=args.profile, region=args.region
    )
    mock_client.describe_scheduled_actions.assert_called_once()


@patch("core.services.asg.create_boto_session")
@patch("builtins.input", side_effect=["1", "test-schedule", "1", "2024-04-10 10:00:00"])
def test_add_autoscaling_schedule_with_input(
    mock_input, mock_create_boto_session, capsys
):
    # Setup the parser with arguments as if they were passed from the command line
    args = setup_parser(add_subparsers, ["asg", "schedule", "add"])
    args.config = setup_config(None)

    mock_session = mock_create_boto_session.return_value
    mock_client = mock_session.client.return_value
    mock_client.describe_auto_scaling_groups.return_value = {
        "AutoScalingGroups": [
            {"AutoScalingGroupName": "test-asg", "DesiredCapacity": 2}
        ]
    }
    mock_client.put_scheduled_update_group_action.return_value = {
        "ResponseMetadata": {"HTTPStatusCode": 200}
    }

    add_autoscaling_schedule(args)
    out, _ = capsys.readouterr()

    # Assertions
    assert "Schedule created successfully" in out
    mock_create_boto_session.assert_called_once_with(
        profile=args.profile, region=args.region
    )
    mock_client.put_scheduled_update_group_action.assert_called_once_with(
        AutoScalingGroupName='test-asg',
        ScheduledActionName="test-schedule",
        MinSize=1,
        StartTime="2024-04-10 10:00:00"
    )

    # Verify the input function was called the expected number of times
    assert mock_input.call_count == 4


@patch("core.services.asg.create_boto_session")
def test_add_autoscaling_schedule_with_args(
    mock_create_boto_session, capsys
):
    # Setup the parser with arguments as if they were passed from the command line
    args = setup_parser(add_subparsers, [
        "asg", "schedule", "add", "test-asg", "test-schedule",
        "--desired", "2", "--min", "1", "--max", "8",
        "--start", "2024-04-10 10:00:00"
    ])
    args.config = setup_config(None)

    # Mock the boto session and the autoscaling client
    mock_session = mock_create_boto_session.return_value
    mock_client = mock_session.client.return_value
    mock_client.put_scheduled_update_group_action.return_value = {
        "ResponseMetadata": {"HTTPStatusCode": 200}
    }

    # Call the function under test
    add_autoscaling_schedule(args)
    out, _ = capsys.readouterr()

    # Assertions
    assert "Schedule created successfully" in out
    mock_create_boto_session.assert_called_once_with(
        profile=args.profile, region=args.region
    )
    mock_client.put_scheduled_update_group_action.assert_called_once_with(
        AutoScalingGroupName="test-asg",
        ScheduledActionName="test-schedule",
        MinSize=1,
        StartTime="2024-04-10 10:00:00",
        DesiredCapacity=2,
        MaxSize=8
    )
