"""
Test case for the common functions module.
"""
from unittest.mock import patch
import pytest
from core import common


def test_print_as_table(capsys):
    """
    Test case for the print_as_table function.

    This test verifies that the print_as_table function correctly prints a list of dicts as a table.

    Args:
        capsys: Pytest fixture to capture stdout and stderr.

    Returns:
        None
    """
    data = [
        {"Name": "instance-1", "Environment": "production", "Status": "running"},
        {"Name": "instance-2", "Environment": "development", "Status": "stopped"},
    ]
    common.print_as_table(data)
    out, _ = capsys.readouterr()

    print("\n" + out)

    assert out == (
        "Name        Environment    Status\n"
        "----------  -------------  --------\n"
        "instance-1  production     running\n"
        "instance-2  development    stopped\n"
    )


@patch('core.common.Session')
@pytest.mark.parametrize("profile, region", [("my-profile", "us-west-2"), (None, "us-west-2"), ("my-profile", None)])
def test_create_boto_session(mock_session, profile, region):
    """
    Test case for the create_boto_session function.

    This test verifies that the create_boto_session function sets up the AWS session correctly.

    Returns:
        None
    """
    session_instance = mock_session.return_value
    session_instance.client.return_value.get_caller_identity.return_value = "dummy_identity"

    session = common.create_boto_session(profile, region)

    if profile and region:
        mock_session.assert_called_once_with(profile_name=profile, region_name=region)
    elif profile:
        mock_session.assert_called_once_with(profile_name=profile)
    elif region:
        mock_session.assert_called_once_with(region_name=region)
    session_instance.client.assert_called_once_with('sts')
    session_instance.client.return_value.get_caller_identity.assert_called_once()

    assert session == session_instance
