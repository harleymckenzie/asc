"""
Test cases for the configdriver module.

This module contains the test cases for the configdriver module.
"""

from unittest.mock import patch, mock_open
import configparser
import pytest
from core.configdriver import setup_config, config_path


@pytest.mark.parametrize(
    "tags_input",
    [(""), ("Environment"), ("Project,Version")],
)
@patch("core.configdriver.input")
@patch("builtins.open", new_callable=mock_open, read_data="")
def test_setup_config_initial_setup(mock_file_open, mock_input, tags_input):
    """
    Test case for the setup_config function with initial setup.

    This test verifies that the setup_config function correctly sets up the
    configuration file when called for the first time.

    Args:
        mock_file_open: Mocked open function.
        mock_input: Mocked input function.
        tags_input: Tags to include in displayed tags.

    Returns:
        None
    """
    mock_input.return_value = tags_input
    config = configparser.ConfigParser()
    expected_tags = tags_input if tags_input else "Name"

    setup_config(config, initial_setup=True)

    # Ensure the configuration file was opened correctly
    mock_file_open.assert_called_once_with(config_path, "w", encoding="utf-8")
    mock_file_handle = mock_file_open()

    written_content = "".join(
        call.args[0] for call in mock_file_handle.write.call_args_list
    )
    
    print("Generated config file:\n" + written_content)
    print("displayed_tags value: " + config.get("asc", "displayed_tags"))
    assert written_content == (
        "[asc]\n" f"displayed_tags = {expected_tags}\n\n"
    ), "The configuration file was not written correctly."

    assert (
        config.get("asc", "displayed_tags") == expected_tags
    ), "The displayed_tags configuration object was not updated correctly."


@pytest.mark.parametrize(
    "existing_tags, tags_input",
    [
        ("Name", "Environment"),
        ("Environment", "Name,Environment"),
        ("Name,Environment", ""),
    ],
)
@patch("core.configdriver.input")
@patch("builtins.open", new_callable=mock_open, read_data="")
def test_setup_config_update(
    mock_file_open, mock_input, existing_tags, tags_input
):
    """
    Test case for the setup_config function when called via the 'configure'
    command.

    This test verifies that the setup_config function correctly updates the
    configuration file when called via the 'configure' command.

    Args:
        mock_file_open: Mocked open function.
        mock_input: Mocked input function.

    Returns:
        None
    """
    mock_input.return_value = tags_input
    config = configparser.ConfigParser()
    config["asc"] = {"displayed_tags": existing_tags}

    setup_config(config, initial_setup=False)

    # Ensure the configuration file was opened correctly
    mock_file_open.assert_called_once_with(config_path, "w", encoding="utf-8")
    mock_file_handle = mock_file_open()

    written_content = "".join(
        call.args[0] for call in mock_file_handle.write.call_args_list
    )

    print("Generated config file:\n" + written_content)
    print("displayed_tags value: " + config.get("asc", "displayed_tags"))
    if tags_input:
        assert written_content == (
            "[asc]\n" f"displayed_tags = {tags_input}\n\n"
        ), "The configuration file was not written correctly."
        assert (
            config.get("asc", "displayed_tags") == tags_input
        ), "The displayed_tags configuration object was not updated correctly."
    else:
        assert written_content == (
            "[asc]\n" f"displayed_tags = {existing_tags}\n\n"
        ), "The configuration file was not written correctly."
        assert (
            config.get("asc", "displayed_tags") == existing_tags
        ), "The displayed_tags configuration object was not updated correctly."
