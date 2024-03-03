# asc
'asc' is a simplified version of the AWS CLI.

The purpose of this project is to provide a simple interface to the AWS CLI.

# AWS Simple CLI (`asc`)

The AWS Simple CLI (`asc`) is a command-line interface for managing AWS services with ease. It supports various AWS services, including EC2, RDS, Autoscaling, Redis, and Systems Manager.

This project was created as a way for me to learn about building CLIs in Python and to provide a simple interface to the AWS CLI. It is not meant to replace the AWS CLI, but to provide a simpler way to interact with AWS services.

## Installation

To install `asc` using Homebrew, run:

```shell
brew tap harleymckenzie/asc
brew install asc
```

## Usage

To use `asc`, run:

```
asc [global options] subcommand [subcommand options] [arguments...]
```

### Subcommands

- `configure`: Configure `asc` settings.
- `asg`: Interact with the Autoscaling service.
- `ec2`: Manage EC2 instances.
- `rds`: Work with RDS databases.
- `redis`: Manage Redis instances.
- `ssm`: Use AWS Systems Manager functionalities.

### Global Options

- `-h, --help`: Show help message and exit.
- `--version`: Show the program's version number and exit.
- `--tags TAGS, -t TAGS`: Comma-separated tags to display in output.
- `--profile [PROFILE], -p [PROFILE]`: Specify the AWS profile to use.
- `--region [REGION]`: Define the AWS region to operate in.
- `-v`: Increase verbosity of output.

### Examples

- List EC2 instances:

```
asc ec2 ls
```

- List Autoscaling groups:

```
asc asg ls
```

- Connect to an instance via Session Manager:

```
asc ssm session i-1234567890abcdef0
```

### Service-Specific Subcommands

#### Autoscaling (`asg`)

- `ls`: List autoscaling groups.
- `schedule`: Manage autoscaling schedules.

##### Schedule Subcommands

- `ls`: List autoscaling schedules.
- `add`: Add a new schedule.
- `rm`: Remove an existing schedule.

#### EC2 Service (`ec2`)

- `ls`: List EC2 instances.

#### RDS Service (`rds`)

- `ls`: List RDS instances.

#### Redis Service (`redis`)

- `ls`: List Redis instances.

#### Systems Manager Service (`ssm`)

- `parameter (param)`: Interact with the Parameter Store.
- `session (connect)`: Connect to a managed instance using Session Manager.

##### Parameter Store Subcommands

- `ls`: List parameters.
- `add`: Add or update a parameter.
- `rm`: Remove a parameter.
