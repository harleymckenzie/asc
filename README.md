# asc
AWS Simple CLI (asc) - A simplified interface for AWS operations

## What is asc?

asc is a CLI tool designed to help me upskill in Go while building a tool that I use daily. The goal
is to simplify common AWS operations in a bash-like way that is easy to remember and use.

This is a personal passion project — I am not a software developer, and am using asc as a way to improve my skills in Go.
Because of this:
- Some features may initially be slightly experimental or rough around the edges.
- I will always do my best to test features thoroughly before including them in a release.

Service features will be added gradually, based either on:
- What I find myself needing the most day-to-day, or
- What I feel like working on at the time.

My aim is to implement functionality in a way that feels natural and efficient, rather than overly rigid.
For example:
- Traversing SSM Parameters in a filesystem-like way (e.g., navigating parameters as folders).
- Creating an ASG scheduled action quickly via CLI — ideally faster and less painful than doing it manually in the AWS console.

## Installation

```sh
brew tap harleymckenzie/asc
brew install asc
```


## Service Implementation

Below is a table of service features and commands that I plan to implement, and their current status.

_**\*** Partly implemented. Missing some features that I hope to add in the future_
| Service        | Command / Subcommand      | Status | Notes / Features                                                            |
|:---------------|:--------------------------|:------:|:----------------------------------------------------------------------------|
| ASG            | ls                        | ✓      | List ASGs                                                                   |
| ASG            | modify                    | ✓*     | Modify ASGs, supports relative and absolute values for min, max and desired capacity<br><sub>_\* Currently supports min, max and desired capacity only_</sub>   |
| ASG            | detach                    | ✗      | The specified instance will be removed from the ASG, but will not be terminated. |
| ASG            | schedule add              | ✓      | Add schedule to ASG, supports human friendly time input                     |
| ASG            | schedule ls               | ✓      | List ASG schedules                                                          |
| ASG            | schedule rm               | ✓      | Remove schedule from ASG                                                    |
| ASG            | show / describe           | ✗      | Show ASG details                                                            |
| CloudFormation | events                    | ✗      | List CloudFormation stack events                                            |
| CloudFormation | ls                        | ✓      | List CloudFormation stacks                                                  |
| CloudFormation | rm                        | ✗      | Delete CloudFormation stacks                                                |
| CloudFormation | show / describe           | ✗      | Show CloudFormation stack details                                           |
| CloudFormation | parameter ls              | ✗      | List CloudFormation stack parameters                                        |
| CloudFormation | parameter edit            | ✗      | Edit CloudFormation stack parameters                                        |
| EC2            | ls                        | ✓      | List EC2 instances                                                          |
| EC2            | modify                    | ✗      | Modify EC2 instances                                                        |
| EC2            | show / describe           | ✗      | Show EC2 instance details                                                   |
| EC2            | start                     | ✓      | Start EC2 instances                                                         |
| EC2            | stop                      | ✓      | Stop EC2 instances                                                          |
| EC2            | restart                   | ✓      | Restart EC2 instances                                                       |
| EC2            | rm / terminate            | ✓      | Terminate EC2 instances                                                     |
| EC2            | ami cp                    | ✗      | Copy EC2 AMI                                                                |
| EC2            | ami ls                    | ✓      | List EC2 AMIs                                                               |
| EC2            | ami rm                    | ✗      | Remove EC2 AMI                                                              |
| EC2            | ami show                  | ✓      | Show EC2 AMI details                                                        |
| EC2            | ami rm                    | ✗      | Remove EC2 AMI                                                              |
| EC2            | security-group add        | ✗      | Add EC2 security group rule                                                 |
| EC2            | security-group ls         | ✓      | List EC2 security groups                                                    |
| EC2            | security-group rm         | ✗      | Remove EC2 security group                                                   |
| EC2            | security-group show       | ✓      | Show EC2 security group details                                             |
| EC2            | security-group rule add   | ✗      | Add EC2 security group rule                                                 |
| EC2            | security-group rule rm    | ✗      | Remove EC2 security group rule                                              |
| EC2            | volume create             | ✗      | Create EC2 volume                                                           |
| EC2            | volume ls                 | ✗      | List EC2 volumes                                                            |
| EC2            | volume show               | ✗      | Show EC2 volume details                                                     |
| EC2            | volume rm                 | ✗      | Remove EC2 volume                                                           |
| EC2            | snapshot ls               | ✗      | List EC2 snapshots                                                          |
| EC2            | snapshot show             | ✗      | Show EC2 snapshot details                                                   |
| EC2            | snapshot rm               | ✗      | Remove EC2 snapshot                                                         |
| ECS            | ls                        | ✗      | List ECS clusters, services, tasks                                          |
| ECS            | modify                    | ✗      | Modify ECS clusters and services                                            |
| ECS            | rm / terminate            | ✗      | Terminate ECS tasks                                                         |
| ECS            | schedule add              | ✗      | Add schedule to ECS services                                                |
| ECS            | schedule ls               | ✗      | List ECS schedules                                                          |
| ECS            | schedule rm               | ✗      | Remove schedule from ECS services                                           |
| ElastiCache    | ls                        | ✓      | List ElastiCache clusters                                                   |
| ElastiCache    | modify                    | ✗      | Modify ElastiCache clusters                                                 |
| ElastiCache    | rm / terminate            | ✗      | Terminate ElastiCache clusters                                              |
| ElastiCache    | show / describe           | ✗      | Show ElastiCache instance details                                           |
| ELB            | ls                        | ✓      | List Elastic Load Balancers                                                 |
| ELB            | modify                    | ✗      | Modify Elastic Load Balancers                                               |
| ELB            | rm                        | ✗      | Terminate Elastic Load Balancers                                            |
| ELB            | show / describe           | ✗      | Show Elastic Load Balancer details                                          |
| ELB            | target-group ls           | ✓      | List Elastic Load Balancer target groups                                    |
| ELB            | target-group add          | ✗      | Add target to Elastic Load Balancer target group                            |
| ELB            | target-group rm           | ✗      | Remove target from Elastic Load Balancer target group                       |
| ELB            | target-group show         | ✗      | Show Elastic Load Balancer target group details                             |
| RDS            | ls                        | ✓      | List RDS instances                                                          |
| RDS            | modify                    | ✗      | Modify RDS instances                                                        |
| RDS            | rm                        | ✗      | Terminate RDS instances                                                     |
| RDS            | show / describe           | ✗      | Show RDS instance details                                                   |
| Route53        | ls                        | ✗      | List Route53 hosted zones and records                                       |
| Route53        | modify                    | ✗      | Modify Route53 hosted zones and records                                     |
| Route53        | rm                        | ✗      | Terminate Route53 hosted zones and records                                  |
| Route53        | show / describe           | ✗      | Show Route53 hosted zone and record details                                 |
| S3             | cp                        | ✗      | Copy S3 objects                                                             |
| S3             | ls                        | ✗      | List S3 buckets and objects                                                 |
| S3             | mv                        | ✗      | Move S3 objects                                                             |
| S3             | rm                        | ✗      | Delete S3 buckets                                                           |
| S3             | show / describe           | ✗      | Show S3 bucket or object details                                            |
| SSM            | document ls               | ✗      | List SSM documents                                                          |
| SSM            | document run              | ✗      | Run SSM documents                                                           |
| SSM            | document rm               | ✗      | Delete SSM documents                                                        |
| SSM            | document show             | ✗      | Show SSM document details                                                   |
| SSM            | parameter add             | ✗      | Add SSM parameters                                                          |
| SSM            | parameter cp              | ✗      | Copy SSM parameters, supports wildcards and cross account copying           |
| SSM            | parameter diff            | ✗      | Diff SSM parameters, supports wildcards and cross account diffing           |
| SSM            | parameter edit            | ✗      | Edit SSM parameters                                                         |
| SSM            | parameter ls              | ✗      | List SSM parameters, supports wildcards                                     |
| SSM            | parameter mv              | ✗      | Move SSM parameters, supports wildcards and cross account moving            |
| SSM            | parameter rm              | ✗      | Delete SSM parameters, supports wildcards                                   |
| SSM            | parameter show            | ✗      | Show SSM parameter details                                                  |
| SSM            | session ls                | ✗      | List SSM sessions                                                           |
| SSM            | session start             | ✗      | Start SSM sessions                                                          |
| SSM            | session stop              | ✗      | Stop SSM sessions                                                           |
| SSM            | session rm                | ✗      | Delete SSM sessions                                                         |
| SSM            | session show              | ✗      | Show SSM session details                                                    |
| VPC            | ls                        | ✗      | List VPCs                                                                   |
| VPC            | modify                    | ✗      | Modify VPCs                                                                 |
| VPC            | rm                        | ✗      | Terminate VPCs                                                              |
| VPC            | show / describe           | ✗      | Show VPC details                                                            |
| VPC            | subnet ls                 | ✗      | List VPC subnets                                                            |
| VPC            | subnet add                | ✗      | Add VPC subnet                                                              |
| VPC            | subnet rm                 | ✗      | Remove VPC subnet                                                           |
| VPC            | subnet show               | ✗      | Show VPC subnet details                                                     |
| VPC            | route-table ls            | ✗      | List VPC route tables                                                       |
| VPC            | route-table add           | ✗      | Add VPC route table                                                         |
| VPC            | route-table rm            | ✗      | Remove VPC route table                                                      |
| VPC            | route-table show          | ✗      | Show VPC route table details                                                |
| VPC            | route-table rule add      | ✗      | Add VPC route table rule                                                    |
| VPC            | route-table rule rm       | ✗      | Remove VPC route table rule                                                 |

### Service Implementation: Other Features
| Description                                                 | Status | Notes / Features                                 |
|:------------------------------------------------------------|:-------|:-------------------------------------------------|
| Shell autocompletion                                        | ✓*     | [Brew Shell Completion](https://docs.brew.sh/Shell-Completion) configuration is required. |
| Customise output fields/columns displayed in tables         | ✗      |                                                  |
| Customise features via configuration file                   | ✗      |                                                  |
| Filesystem-like navigation                                  | ✗      |                                                  |
| Optional terminal UI                                        | ✗      |                                                  |
| Export data to CSV, JSON, or other formats                  | ✗      |                                                  |
| Service agnostic `show` command                             | ✗      |                                                  |
| AWS Profile management                                      | ✗      |                                                  |
| 'Select' resources to avoid repeating identifiers           | ✗      |                                                  |
| Display pricing information on supported resources          | ✗      |                                                  |


## Output Format

By default, most `asc` commands output results in a **table format** for easier readability.
Many commands also support a `--list` flag to produce a **simpler list-style output** if preferred.

Example:
```sh
asc ec2 ls
```
_(Outputs EC2 instances in a table.)_

```sh
asc ec2 ls -l
```
_(Outputs EC2 instances in a basic list format.)_

### Example Output

Example output from listing RDS clusters and instances:

```
╭───────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ RDS Clusters and Instances                                                                                        │
├────────────────────┬─────────────────────────────────────────┬───────────┬──────────────┬────────────────┬────────┤
│ CLUSTER IDENTIFIER │ IDENTIFIER                              │ STATUS    │ ENGINE       │ SIZE           │ ROLE   │
├────────────────────┼─────────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
│ prod-aurora        │ prod-aurora-eu1a                        │ available │ aurora-mysql │ db.r6g.2xlarge │ Writer │
│                    ├─────────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
│                    │ prod-aurora-eu1c                        │ available │ aurora-mysql │ db.r6g.2xlarge │ Reader │
│                    ├─────────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
│                    │ reporting-aurora                        │ available │ aurora-mysql │ db.t4g.large   │ Reader │
├────────────────────┼─────────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
│ testing-cluster    │ aurora-testing                          │ available │ aurora-mysql │ db.t3.medium   │ Writer │
│                    ├─────────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
│                    │ aurora-legacy-cluster                   │ available │ aurora-mysql │ db.t4g.medium  │ Writer │
├────────────────────┼─────────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
│ legacy-cluster     │ legacy-reporting-aurora                 │ available │ aurora-mysql │ db.t3.medium   │ Reader │
│                    ├─────────────────────────────────────────┼───────────┼──────────────┼────────────────┼────────┤
│                    │ legacy-upgrade-dry-run-cluster          │ available │ aurora-mysql │ db.t4g.medium  │ Writer │
╰────────────────────┴─────────────────────────────────────────┴───────────┴──────────────┴────────────────┴────────╯
```

## Examples

### EC2

#### List all EC2 instances
```sh
asc ec2 ls
```

#### List all EC2 instances showing AMI ID and private IP
```sh
asc ec2 ls -A -P
```

#### List all EC2 instances sorted by launch time
```sh
asc ec2 ls -t
```

#### Output EC2 instances in a simple list format
```sh
asc ec2 ls -l
```

### Auto Scaling Groups (ASG)

#### List all Auto Scaling Groups
```sh
asc asg ls
```

#### List all Auto Scaling Groups showing ARNs
```sh
asc asg ls --arn
```

#### List all Auto Scaling Groups sorted by number of instances
```sh
asc asg ls -i
```

#### List instances in a specific Auto Scaling Group
```sh
asc asg ls my-asg-name
```

### ASG Schedules

#### List all schedules across all Auto Scaling Groups
```sh
asc asg ls schedules
```

#### List schedules for a specific Auto Scaling Group
```sh
asc asg ls schedules my-asg-name
```

#### Add a schedule to an Auto Scaling Group with minimum and maximum size set, at 10:00am on 25th April 2025
```sh
asc asg schedule add my-schedule -a my-asg -m 4 -M 8 -s 'Friday 10:00'
```

#### Add a schedule to an Auto Scaling Group with desired capacity set, at 10:00am on 25th April 2025
```sh
asc asg schedule add my-schedule -a my-asg -d 8 -s '10:00am 25/04/2025'
```

