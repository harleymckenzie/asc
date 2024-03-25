"""
Test case for the Elasticache Redis module.

This module contains the test cases for the Elasticache Redis module.
"""
from unittest.mock import patch
import pytest
from core.services import redis
from .test_utils import setup_parser, setup_config


@pytest.mark.parametrize(
    "arg_list, displayed_tags, expected_output",
    [
        (
            ['redis', 'ls', '--profile', 'test-profile', '--region', 'eu-west-1', '--endpoint'],
            None,
            "Cluster Id       Type             Status     Endpoint\n"
            "---------------  ---------------  ---------  ---------------------------------------------\n"
            "redis-cluster-1  cache.t3.micro   available  redis-cluster-1.0001.use1.cache.amazonaws.com\n"
            "redis-cluster-2  cache.m7g.large  available  redis-cluster-2.0001.use1.cache.amazonaws.com\n"
            "redis-test-2     cache.t2.micro   available  redis-test-2.0001.use1.cache.amazonaws.com\n"
        ),
        (
            ['redis', 'ls', '--profile', 'test-profile', '--region', 'eu-west-1'],
            "Environment",
            "Cluster Id       Type             Status     Environment\n"
            "---------------  ---------------  ---------  -------------\n"
            "redis-cluster-1  cache.t3.micro   available  production\n"
            "redis-cluster-2  cache.m7g.large  available  production\n"
            "redis-test-2     cache.t2.micro   available  production\n"
        ),
        (
            ['redis', 'ls', '--profile', 'test-profile', '--region', 'eu-west-1', '--sort-by', 'Type'],
            None,
            "Cluster Id       Type             Status\n"
            "---------------  ---------------  ---------\n"
            "redis-cluster-2  cache.m7g.large  available\n"
            "redis-test-2     cache.t2.micro   available\n"
            "redis-cluster-1  cache.t3.micro   available\n"
        )
    ],
    ids=[
        "Output endpoint, No tags",
        "tags: Environment",
        "Sort by Type, No tags"
    ]
)
@patch('core.services.redis.create_boto_session')
def test_list_redis_instances(mock_create_boto_session, arg_list,
                              displayed_tags, expected_output, capsys):
    """
    Test case for the list_redis_instances function.

    This test verifies that the list_redis_instances function correctly retrieves
    and prints the details of Redis instances.

    Args:
        mock_create_boto_session: Mocked function to create a Boto session.
        capsys: Pytest fixture to capture stdout and stderr.

    Returns:
        None
    """
    args = setup_parser(redis.add_subparsers, arg_list)
    args.config = setup_config(displayed_tags)

    mock_session = mock_create_boto_session.return_value
    mock_client = mock_session.client.return_value
    mock_client.describe_cache_clusters.return_value = {
        "CacheClusters": [
            {
                "CacheClusterId": "redis-cluster-1",
                "CacheNodeType": "cache.t3.micro",
                "ARN": "arn:aws:elasticache:eu-west-1:123456789012:cluster:redis-cluster-1",
                "Engine": "redis",
                "EngineVersion": "5.0.6",
                "CacheClusterStatus": "available",
                "CacheNodes": [
                    {
                        "CacheNodeId": "0001",
                        "CacheNodeStatus": "available",
                        "Endpoint": {
                            "Address": "redis-cluster-1.0001.use1.cache.amazonaws.com",
                            "Port": 6379
                        }
                    }
                ]
            },
            {
                "CacheClusterId": "redis-test-2",
                "CacheNodeType": "cache.t2.micro",
                "ARN": "arn:aws:elasticache:eu-west-1:123456789012:cluster:redis-test-2",
                "Engine": "redis",
                "EngineVersion": "7.0.5",
                "CacheClusterStatus": "available",
                "CacheNodes": [
                    {
                        "CacheNodeId": "0001",
                        "CacheNodeStatus": "available",
                        "Endpoint": {
                            "Address": "redis-test-2.0001.use1.cache.amazonaws.com",
                            "Port": 6379
                        }
                    }
                ]
            },
            {
                "CacheClusterId": "redis-cluster-2",
                "CacheNodeType": "cache.m7g.large",
                "ARN": "arn:aws:elasticache:eu-west-1:123456789012:cluster:redis-cluster-2",
                "Engine": "redis",
                "EngineVersion": "7.0.5",
                "CacheClusterStatus": "available",
                "CacheNodes": [
                    {
                        "CacheNodeId": "0001",
                        "CacheNodeStatus": "available",
                        "Endpoint": {
                            "Address": "redis-cluster-2.0001.use1.cache.amazonaws.com",
                            "Port": 6379
                        }
                    }
                ]
            }
        ]
    }
    mock_client.list_tags_for_resource.return_value = {
        "TagList": [
            {"Key": "Environment", "Value": "production"},
            {"Key": "Owner", "Value": "John Doe"}
        ]
    }

    redis.list_redis_instances(args)
    out, _ = capsys.readouterr()
    print("\n" + out)
    assert out.strip() == expected_output.strip()

    # Verify that external calls were made as expected
    mock_create_boto_session.assert_called_once_with(profile=args.profile,
                                                     region=args.region)
    mock_client.describe_cache_clusters.assert_called_once()
    # Add any other necessary assertions for your mock calls
