#!/usr/bin/env python
import subprocess
from packaging.version import Version


def get_last_tag():
    # Retrieve the last semver tag from the git history
    try:
        last_tag = subprocess.check_output(["git", "tag", "--sort=-v:refname", "--sort=-creatordate"]).decode().split('\n')[0]
    except subprocess.CalledProcessError:
        last_tag = "0.0.0"
    return last_tag


def get_commit_messages_since_last_tag(last_tag):
    # Retrieve commit messages since the last tag
    commit_messages = subprocess.check_output(["git", "log", f"{last_tag}..HEAD", "--pretty=format:%s"]).decode().split('\n')
    return commit_messages


def determine_bump_level(commit_messages, default_bump="minor"):
    # Initialize bump levels found
    found_levels = {"#major": False, "#minor": False, "#patch": False}
    
    # Determine version bump based on commit messages
    for message in commit_messages:
        if "#major" in message:
            found_levels["#major"] = True
        elif "#minor" in message:
            found_levels["#minor"] = True
        elif "#patch" in message:
            found_levels["#patch"] = True
        elif "#none" in message:
            return "none"
    
    # Determine the bump level based on what was found
    if found_levels["#major"]:
        return "major"
    elif found_levels["#minor"]:
        return "minor"
    elif found_levels["#patch"]:
        return "patch"
    else:
        return default_bump


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
    print(new_version)


if __name__ == "__main__":
    main()
