# asc
AWS Simple CLI (asc) - A simplified interface for AWS operations

# What is asc?

asc is a CLI tool that allows you to interact with AWS services in a simplified way. It is designed to be easy to use and understand, and to provide a consistent interface for interacting with AWS services.

## Installation

```sh
brew tap harleymckenzie/asc
brew install asc
```

# Subcommands
- `ec2` - EC2 operations
- `rds` - RDS operations
- `elasticache` - ElastiCache operations
- `asg` - ASG operations

# Examples

List all EC2 instances, diplaying and sorting by time created:

```sh
asc ec2 ls -Lt
```

List all RDS resources, sorted by instance identifier:

```sh
asc rds ls -n
```

List all ElastiCache clusters, sorted by type:

```sh
asc elasticache ls -T
```

List all ASGs, sorted by name:

```sh
asc asg ls -n
```

List all instances in an ASG:

```sh
asc asg ls <asg-name>
```

