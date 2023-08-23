from unittest.mock import patch
from services.common import print_as_table


def test_print_as_table():
    # Sample data to test
    items = [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]

    # Expected output string
    expected_out = "  id  name\n" "----  ------\n" "   1  Alice\n" "   2  Bob"

    # Print a newline directly to format pytest output
    print()

    # Mock the print function and capture its arguments
    with patch("builtins.print", side_effect=print) as mocked_print:
        print_as_table(items)

    # Check if the mock was called with the expected output
    mocked_print.assert_called_once_with(expected_out)
