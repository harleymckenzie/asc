"""
Unit tests for the SSM service module.

Functions:
- test_list_ssm_parameters: Test the list_ssm_parameters function.
- test_add_or_update_parameter: Test the add_or_update_parameter function.
- test_remove_parameter: Test the remove_parameter function.
"""
import pytest
from asc.services import ssm
from unittest import mock
from .test_utils import setup_args, run_and_capture_output, mock_ssm_client


def create_parameters(mock_ssm_client):
    """
    Create parameters for testing.
    """
    mock_ssm_client.put_parameter(
        Name="test-parameter",
        Value="test-value",
        Type="String",
        Overwrite=True
    ),
    mock_ssm_client.put_parameter(
        Name="test-secure-parameter",
        Value="test-secure-value",
        Type="SecureString",
        Overwrite=True
    )


@pytest.mark.parametrize("values, decrypt", [(False, False), (True, False), (True, True)])
def test_list_ssm_parameters(mock_ssm_client, values, decrypt):
    """
    Test the list_ssm_parameters function.
    If 'values' is True, String type parameters should display their values,
    and SecureString type parameters should display their values as '*****'.
    If 'decrypt' is True, SecureString type parameters should display their decrypted values.
    """
    args = setup_args()
    args.values, args.decrypt = values, decrypt
    
    create_parameters(mock_ssm_client)
    output = run_and_capture_output(ssm.list_ssm_parameters, args)

    # Print output for debugging
    print("\n" + output[0])

    # Check for expected headers
    expected_headers = ["Name", "Type", "Last Modified"]
    assert all(header in output[0] for header in expected_headers)

    # Confirm that the parameters are displayed as expected
    # "Values" column should not be displayed if 'values' is False
    # 
    if values:
        assert "test-value" in output[0]
        assert "*****" in output[0] if not decrypt else "test-secure-value" in output[0]
    else:
        assert "test-value" not in output[0]
        assert "*****" not in output[0]
        assert "test-secure-value" not in output[0]
