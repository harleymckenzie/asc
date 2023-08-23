import pytz
from .common import print_as_table


def add_subparsers(subparsers, global_parser) -> None:
    asg_parser = subparsers.add_parser(
        "asg",
        help="Autoscaling service",
        description="Autoscaling service",
        epilog="""Example: asc asg ls""",
        parents=[global_parser],
    )
    asg_parser.set_defaults(func=lambda args: asg_parser.print_help())
    asg_subparsers = asg_parser.add_subparsers(
      help="Description:",
      dest="subcommand"
    )

    # ASG specific subcommands
    asg_list_parser = asg_subparsers.add_parser(
        "ls",
        help="List autoscaling groups",
        description="List autoscaling groups",
        epilog="""Example: asc asg ls""",
        parents=[global_parser],
    )
    asg_list_parser.set_defaults(func=list_autoscaling_groups)

    # ASG schedule subcommands
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
        help="Description:", dest="subcommand"
    )

    # ASG schedule list subcommand
    schedule_list_parser = schedule_subparsers.add_parser(
        "ls",
        help="List autoscaling schedules",
        description="List autoscaling schedules",
        epilog="""Example: asc asg schedule ls""",
        parents=[global_parser],
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

    schedule_add_parser.add_argument("--asg", help="Name of the ASG")
    schedule_add_parser.add_argument("--name", help="Name of the schedule")
    schedule_add_parser.add_argument("--min", help="Min size of the ASG")
    schedule_add_parser.add_argument("--start", help="Schedule start time")
    schedule_add_parser.set_defaults(func=add_autoscaling_schedule)

    # ASG schedule rm subcommand
    schedule_rm_parser = schedule_subparsers.add_parser(
        "rm",
        help="Remove autoscaling schedule",
        description="Remove autoscaling schedule",
        epilog="""Example: asc asg schedule rm my-schedule my-asg""",
        parents=[global_parser],
    )
    schedule_rm_parser.add_argument("--name", help="Name of the schedule")
    schedule_rm_parser.add_argument("--asg", help="Name of the ASG")
    schedule_rm_parser.set_defaults(func=rm_autoscaling_schedule)


def list_autoscaling_groups(args):
    """
    List ASGs
    """
    instance_list = []
    asg = args.session.client("autoscaling")
    response = asg.describe_auto_scaling_groups()

    for asg in response["AutoScalingGroups"]:
        instance = {}
        instance["Name"] = asg["AutoScalingGroupName"]
        instance["Min"] = asg["MinSize"]
        instance["Max"] = asg["MaxSize"]
        instance["Desired"] = asg["DesiredCapacity"]
        instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i["Name"])
    print_as_table(instances)


def list_autoscaling_schedules(args):
    """
    List ASG Schedules
    """
    instance_list = []
    asg = args.session.client("autoscaling")
    response = asg.describe_scheduled_actions()

    for schedule in response["ScheduledUpdateGroupActions"]:
        asg_response = asg.describe_auto_scaling_groups(
            AutoScalingGroupNames=[schedule["AutoScalingGroupName"]]
        )

        # Store the name of the ASG for each scheduled action
        for tag in asg_response["AutoScalingGroups"][0]["Tags"]:
            if tag["Key"] == "Name":
                asg_name = tag["Value"]

        # Create a dict for each scheduled action
        instance = {}
        instance["Name"] = schedule["ScheduledActionName"]
        instance["ASG"] = asg_name
        instance["Start Time"] = (
            schedule["StartTime"]
            .astimezone(pytz.timezone("Europe/London"))
            .strftime("%Y-%m-%d %H:%M:%S %Z")
        )
        if "Recurrence" in schedule:
            instance["Recurrence"] = schedule["Recurrence"]
        if "DesiredCapacity" in schedule:
            instance["Desired"] = schedule["DesiredCapacity"]
        if "MinSize" in schedule:
            instance["Min"] = schedule["MinSize"]
        if "MaxSize" in schedule:
            instance["Max"] = schedule["MaxSize"]
        instance_list.append(instance)

    instances = sorted(instance_list, key=lambda i: i["Start Time"])
    print_as_table(instances)


def add_autoscaling_schedule(args):
    """
    Add ASG Schedule
    """
    asg = args.session.client("autoscaling")

    print("Available ASGs:")
    asg_response = asg.describe_auto_scaling_groups()
    for asg_group in asg_response["AutoScalingGroups"]:
        for tag in asg_group["Tags"]:
            if tag["Key"] == "Name":
                print(tag["Value"], "\n")

    request_params = {}

    # Get the parameters from user input if not provided
    asg_name = args.asg or input("ASG Name: ")
    schedule_name = args.name or input("Schedule Name: ")
    min_size = args.min or input("Min Size: ")
    start_time = args.start or input("Start Time (YYYY-MM-DD HH:MM:SS): ")
    max_size = input("Max Size (optional): ")
    desired_size = input("Desired Size (optional): ")
    recurrence = input("Recurrence (10 0 * * *) (optional): ")

    # Get the AutoScalingGroupName from the Name tag
    asg_response = asg.describe_auto_scaling_groups()
    for asg_group in asg_response["AutoScalingGroups"]:
        for tag in asg_group["Tags"]:
            if tag["Key"] == "Name" and tag["Value"] == asg_name:
                asg_name = asg_group["AutoScalingGroupName"]

    # Add the parameters to the request_params dict
    request_params["AutoScalingGroupName"] = asg_name
    request_params["ScheduledActionName"] = schedule_name
    request_params["StartTime"] = start_time
    request_params["MinSize"] = int(min_size)
    if max_size:
        request_params["MaxSize"] = int(max_size)
    if desired_size:
        request_params["DesiredCapacity"] = int(desired_size)
    if recurrence:
        request_params["Recurrence"] = recurrence

    # Create the scheduled action
    response = asg.put_scheduled_update_group_action(**request_params)

    if response["ResponseMetadata"]["HTTPStatusCode"] == 200:
        print("Schedule created successfully")
    else:
        print("Error creating schedule")


def rm_autoscaling_schedule(args):
    """
    Remove ASG Schedule
    """
    asg = args.session.client("autoscaling")

    # Get the parameters from user input if not provided
    asg_name = args.asg if args.asg else input("ASG Name: ")
    schedule_name = args.name if args.name else input("Schedule Name: ")

    # Get the AutoScalingGroupName from the Name tag
    asg_response = asg.describe_auto_scaling_groups()
    for asg_group in asg_response["AutoScalingGroups"]:
        for tag in asg_group["Tags"]:
            if tag["Key"] == "Name" and tag["Value"] == asg_name:
                asg_name = asg_group["AutoScalingGroupName"]

    # Remove the scheduled action
    response = asg.delete_scheduled_action(
        AutoScalingGroupName=asg_name, ScheduledActionName=schedule_name
    )

    if response["ResponseMetadata"]["HTTPStatusCode"] == 200:
        print("Schedule removed successfully")
    else:
        print("Error removing schedule")
