"""
Test case for the Elasticache Redis module.

This module contains the test cases for the Elasticache Redis module.
"""
from unittest.mock import patch
import pytest
from core.services import redis
from .test_utils import setup_args


@patch('core.services.redis.create_boto_session')
@pytest.mark.parametrize("displayed_tags, endpoint", [("Name", True), ("Name,Environment", False), (None, False)])
def test_list_redis_instances(create_boto_session, displayed_tags, endpoint, capsys):
    """
    Test case for the list_redis_instances function.

    This test verifies that the list_redis_instances function correctly retrieves
    and prints the details of Redis instances.

    Args:
        create_boto_session: Mocked function to create a Boto session.
        displayed_tags: Tags to display.
        capsys: Pytest fixture to capture stdout and stderr.

    Returns:
        None
    """
    args = setup_args(displayed_tags=displayed_tags)
    args.endpoint = endpoint

    mock_session = create_boto_session.return_value
    mock_client = mock_session.client.return_value

    mock_client.describe_cache_clusters.return_value = {
        "CacheClusters": [
            {
                "CacheClusterId": "redis-cluster-1",
                "CacheNodeType": "cache.t2.micro",
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
            }
        ]
    }
    mock_client.list_tags_for_resource.return_value = {
        "TagList": [
            {"Key": "Name", "Value": "production-redis-cluster"},
            {"Key": "Environment", "Value": "production"},
            {"Key": "Owner", "Value": "John Doe"}
        ]
    }

    redis.list_redis_instances(args)
    out, _ = capsys.readouterr()

    print("\n" + out)

    if displayed_tags == "Name" and endpoint:
        assert out == (
            "Name                      Cluster Id       Type            Status     Endpoint\n"
            "------------------------  ---------------  --------------  ---------  ---------------------------------------------\n"
            "production-redis-cluster  redis-cluster-1  cache.t2.micro  available  redis-cluster-1.0001.use1.cache.amazonaws.com\n"
        )
    elif displayed_tags == "Name,Environment" and not endpoint:
        assert out == (
            "Name                      Cluster Id       Type            Status     Environment\n"
            "------------------------  ---------------  --------------  ---------  -------------\n"
            "production-redis-cluster  redis-cluster-1  cache.t2.micro  available  production\n"
        )
    else:
        assert out == (
            "Cluster Id       Type            Status\n"
            "---------------  --------------  ---------\n"
            "redis-cluster-1  cache.t2.micro  available\n"
        )
