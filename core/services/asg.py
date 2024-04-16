"""
ASG service.

This module contains functions to interact with the Autoscaling service.

It provides functionality to list autoscaling groups, list autoscaling
schedules, add autoscaling schedules and remove autoscaling schedules.

Functions:
- add_subparsers(subparsers, global_parser) -> None
- list_autoscaling_groups(args)
- list_autoscaling_schedules(args)
- add_autoscaling_schedule(args)
"""

import logging
from typing import Any
from ..common import (
    subparser_register,
    create_boto_session,
    print_as_table,
    apply_tags,
)


logger = logging.getLogger(__name__)


@subparser_register("asg")
def add_subparsers(subparsers: Any, global_parser: Any) -> None:
    """
    Adds subparsers to the given subparsers object for the autoscaling service.

    Args:
        subparsers: The subparsers object to add the subcommands to.
        global_parser: The global parser object to inherit options from.

    Returns:
        None
    """
    asg_parser = subparsers.add_parser(
        "asg",
        help="Autoscaling service",
        description="Autoscaling service",
        epilog="""Example: asc asg ls""",
        parents=[global_parser],
    )
    asg_parser.set_defaults(func=lambda args: asg_parser.print_help())
    asg_subparsers = asg_parser.add_subparsers(
        help="", metavar="subcommand", dest="subcommand"
    )

    # ASG specific subcommands
    asg_list_parser = asg_subparsers.add_parser(
        "ls",
        help="List autoscaling groups",
        description="List autoscaling groups",
        epilog="""Example: asc asg ls""",
        parents=[global_parser],
    )
    asg_list_parser.add_argument(
        "--sort-by",
        help="The field to sort by",
        default="ASG Name",
    )
    asg_list_parser.add_argument(
        "--sort-order",
        help="Specify sort order: 'asc' for ascending or 'desc' for descending",
        default="asc",
    )
    asg_list_parser.set_defaults(func=list_autoscaling_groups)

    # ASG Scheduler Subparser
    add_scheduler_parser(asg_subparsers, global_parser)


def add_scheduler_parser(asg_subparsers: Any, global_parser: Any):
    """
    Adds subparsers for ASG schedule commands.
    """
    schedule_parser = asg_subparsers.add_parser(
        "schedule",
        help="Autoscaling schedule subcommands",
        description="Autoscaling schedule subcommands",
        epilog="""Example: asc asg schedule ls""",
        parents=[global_parser],
    )
    schedule_parser.set_defaults(
        func=lambda args: schedule_parser.print_help()
    )
    schedule_subparsers = schedule_parser.add_subparsers(
        help="", metavar="subcommand", dest="subcommand"
    )

    # ASG schedule list subcommand
    schedule_list_parser = schedule_subparsers.add_parser(
        "ls",
        help="List autoscaling schedules",
        description="List autoscaling schedules",
        epilog="""Example: asc asg schedule ls""",
        parents=[global_parser],
    )
    schedule_list_parser.add_argument(
        "--sort-by",
        help="The field to sort by",
        default="ASG Name",
    )
    schedule_list_parser.add_argument(
        "--sort-order",
        help="Specify sort order: 'asc' for ascending or 'desc' for descending",
        default="asc",
    )
    schedule_list_parser.set_defaults(func=list_autoscaling_schedules)

    # ASG schedule add subcommand
    schedule_add_parser = schedule_subparsers.add_parser(
        "add",
        help="Add autoscaling schedule",
        description="Add autoscaling schedule",
        epilog="""Example: asc asg schedule add --asg my-asg
                --name my-schedule --min 1 --start 2017-01-01T00:00:00Z""",
        parents=[global_parser],
    )

    schedule_add_parser.add_argument(
        "asg_name", help="Name of the ASG", default=None, nargs="?"
    )
    schedule_add_parser.add_argument(
        "schedule_name",
        help="Name of the schedule",
        default=None,
        nargs="?",
    )
    schedule_add_parser.add_argument(
        "--desired", help="Desired capacity of the ASG", type=int, default=None
    )
    schedule_add_parser.add_argument(
        "--min", help="Min size of the ASG", type=int, default=None
    )
    schedule_add_parser.add_argument(
        "--max", help="Max size of the ASG", type=int, default=None
    )
    schedule_add_parser.add_argument(
        "--start", help="Schedule start time", default=None
    )
    schedule_add_parser.set_defaults(func=add_autoscaling_schedule)

    # ASG schedule rm subcommand
    schedule_rm_parser = schedule_subparsers.add_parser(
        "rm",
        help="Remove autoscaling schedule",
        description="Remove autoscaling schedule",
        epilog="""Example: asc asg schedule rm my-schedule my-asg""",
        parents=[global_parser],
    )
    schedule_rm_parser.add_argument(
        "asg_name", help="Name of the ASG", nargs="?"
    )
    schedule_rm_parser.add_argument(
        "schedule_name", help="Name of the schedule", nargs="?"
    )
    schedule_rm_parser.set_defaults(func=rm_autoscaling_schedule)


def list_autoscaling_groups(args: Any) -> None:
    """
    List all autoscaling groups.

    Args:
        args: Arguments passed by the user.

    Prints:
        A table displaying the details of all autoscaling groups.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    asg_client = session.client("autoscaling")
    displayed_tags_list = args.config.get(
        "asc", "displayed_tags", fallback=""
    ).split(",")
    instance_list = []

    try:
        response = asg_client.describe_auto_scaling_groups()
    except Exception as e:
        logger.error("Failed to list Auto Scaling Groups: %s", e)
        exit(1)

    for instance_data in response["AutoScalingGroups"]:
        instance = {
            "ASG Name": instance_data["AutoScalingGroupName"],
            "Min": instance_data["MinSize"],
            "Max": instance_data["MaxSize"],
            "Desired": instance_data["DesiredCapacity"],
        }

        instance = apply_tags(instance, instance_data, displayed_tags_list)
        instance_list.append(instance)

    key_order = [
        "ASG Name",
        "Min",
        "Max",
        "Desired",
    ] + displayed_tags_list
    print_as_table(
        instance_list,
        key_order=key_order,
        sort_key=args.sort_by,
        sort_order=args.sort_order,
    )


def list_autoscaling_schedules(args: Any) -> None:
    """
    List all autoscaling schedules.

    Args:
        args: Arguments passed by the user.

    Prints:
        A table displaying the details of all autoscaling schedules.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    asg_client = session.client("autoscaling")
    instance_list = []

    try:
        response = asg_client.describe_scheduled_actions()
    except Exception as e:
        logger.error("Failed to list Auto Scaling Groups: %s", e)
        exit(1)

    for instance_data in response["ScheduledUpdateGroupActions"]:
        instance = {
            "ASG Name": instance_data["AutoScalingGroupName"],
            "Name": instance_data["ScheduledActionName"],
            "Start Time": instance_data["StartTime"],
        }

        # Only include the following fields if they're present
        for key, new_key in [
            ("Recurrence", "Recurrence"),
            ("DesiredCapacity", "Desired"),
            ("MinSize", "Min"),
            ("MaxSize", "Max"),
        ]:
            if key in instance_data:
                instance[new_key] = instance_data[key]

        instance_list.append(instance)

    key_order = [
        "ASG Name",
        "Name",
        "Start Time",
        "Min",
        "Max",
        "Desired",
    ]
    print_as_table(
        instance_list,
        key_order=key_order,
        sort_key=args.sort_by,
        sort_order=args.sort_order,
    )


def add_autoscaling_schedule(args: Any) -> None:
    """
    Add a new autoscaling schedule.

    Args:
        args: Arguments including ASG Name, Schedule Name, Min Size, etc.

    Prints:
        Confirmation message upon successful creation, or error if failure.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    asg_client = session.client("autoscaling")

    # Get the parameters from user input if not provided
    schedule_parameters = {
        "AutoScalingGroupName": args.asg_name or select_asg(session),
        "ScheduledActionName": args.schedule_name or input("Schedule Name: "),
        "MinSize": args.min or int(input("Min Size: ")),
        "StartTime": args.start or input("Start Time (YYYY-MM-DD HH:MM:SS): "),
    }

    if args.desired:
        schedule_parameters["DesiredCapacity"] = int(args.desired)
    if args.max:
        schedule_parameters["MaxSize"] = int(args.max)

    try:
        response = asg_client.put_scheduled_update_group_action(
            **schedule_parameters
        )
    except Exception as e:
        logger.error("Failed to create schedule: %s", e)
        exit(1)

    if response["ResponseMetadata"]["HTTPStatusCode"] == 200:
        print("Schedule created successfully")
    else:
        print("Error creating schedule")
        exit(1)


def select_asg(session) -> str:
    """
    Display a list of ASGs and let the user select one.

    Returns:
        The name of the selected ASG.
    """
    asg_client = session.client("autoscaling")
    response = asg_client.describe_auto_scaling_groups()
    asgs = [
        asg["AutoScalingGroupName"] for asg in response["AutoScalingGroups"]
    ]
    for i, asg in enumerate(asgs, start=1):
        print(f"{i}. {asg}")
    selection = int(input("Select an ASG by number: ")) - 1
    return asgs[selection]


def rm_autoscaling_schedule(args: Any) -> None:
    """
    Remove an autoscaling schedule.

    Args:
        args: Arguments including ASG Name and Schedule Name.

    Prints:
        Confirmation message upon successful removal, or error if failure.
    """
    session = create_boto_session(profile=args.profile, region=args.region)
    asg_client = session.client("autoscaling")

    # Get the parameters from user input if not provided
    asg_name = args.asg_name if args.asg_name else input("ASG Name: ")
    schedule_name = (
        args.schedule_name if args.schedule_name else input("Schedule Name: ")
    )

    # Get the AutoScalingGroupName from the Name tag
    # asg_response = asg_client.describe_auto_scaling_groups()
    # logger.debug("ASG response: %s", asg_response)
    # for asg_group in asg_response["AutoScalingGroups"]:
    #     for tag in asg_group["Tags"]:
    #         if tag["Key"] == "Name" and tag["Value"] == asg_name:
    #             asg_name = asg_group["AutoScalingGroupName"]
    #             logger.info(asg_name)

    # Remove the scheduled action
    try:
        response = asg_client.delete_scheduled_action(
            AutoScalingGroupName=asg_name, ScheduledActionName=schedule_name
        )
    except Exception as e:
        logger.error("Failed to remove schedule: %s", e)
        exit(1)

    if response["ResponseMetadata"]["HTTPStatusCode"] == 200:
        print("Schedule removed successfully")
    else:
        logger.error("Error removing schedule")
