# asc
'asc' is a simplified version of the AWS CLI.

The purpose of this project is to provide a simple interface to the AWS CLI.

## Usage
```
usage: asc [-h] [--profile [PROFILE]] [--region [REGION]] subcommand ...

positional arguments:
  subcommand           description
    configure          Configure asc
    ec2                EC2 service
    rds                RDS service
    asg                Autoscaling service
    redis              Redis service
    
options:
  -h, --help           show this help message and exit
  --profile [PROFILE]  AWS profile to use
  --region [REGION]    AWS region to use
```

## Examples
### EC2
List EC2 instances
```shell
asc ec2 ls
```

### RDS
List RDS instances
```shell
asc rds ls
```

### Autoscaling
List Autoscaling groups
```shell
asc asg ls
```

List Autoscaling scheduled actions
```shell
asc asg schedule ls
```

Add Autoscaling scheduled action (parameters are optional)
```shell
asc asg schedule add --asg my-asg --name my-schedule --min 1 --start 2017-01-01T00:00:00Z```
```

Remove Autoscaling scheduled action (as above)
```shell
asc asg schedule rm --asg my-asg --name my-schedule
```

### Redis
List Redis clusters
```shell
asc redis ls
```

Tag Redis cluster
```shell
asc redis tag --cluster my-cluster --key my-key --value my-value
```