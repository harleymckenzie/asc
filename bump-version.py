#!/usr/bin/env python
import re
import subprocess
from packaging.version import Version


def get_last_tag():
    # Retrieve the last semver tag from the git history
    try:
        last_tag = subprocess.check_output(["git", "describe", "--tags", "--abbrev=0"]).decode().strip()
    except subprocess.CalledProcessError:
        last_tag = "0.0.0"
    return last_tag


def get_commit_messages_since_last_tag(last_tag):
    # Retrieve commit messages since the last tag
    commit_messages = subprocess.check_output(["git", "log", f"{last_tag}..HEAD", "--pretty=format:%s"]).decode().split('\n')
    return commit_messages


def determine_bump_level(commit_messages, default_bump="minor"):
    # Determine version bump based on commit messages
    bump_level = default_bump
    for message in commit_messages:
        if "#major" in message:
            return "major"
        elif "#minor" in message and bump_level != "major":
            bump_level = "minor"
        elif "#patch" in message and bump_level not in ["major", "minor"]:
            bump_level = "patch"
        elif "#none" in message:
            return "none"
    return bump_level


def calculate_new_version(last_tag, bump_level):
    # Calculate the new version based on the last tag and the determined bump level
    if bump_level == "none":
        return last_tag
    version = Version(last_tag)
    if bump_level == "major":
        new_version = Version(f"{version.major + 1}.0.0")
    elif bump_level == "minor":
        new_version = Version(f"{version.major}.{version.minor + 1}.0")
    elif bump_level == "patch":
        new_version = Version(f"{version.major}.{version.minor}.{version.micro + 1}")
    return str(new_version)


def main():
    last_tag = get_last_tag()
    commit_messages = get_commit_messages_since_last_tag(last_tag)
    bump_level = determine_bump_level(commit_messages)
    new_version = calculate_new_version(last_tag, bump_level)
    print(f"New version: {new_version}")


if __name__ == "__main__":
    main()
