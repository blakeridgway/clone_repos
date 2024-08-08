#!/bin/bash

# Prompt for Github username

read -p "Enter your Github username: " GITHUB_USER

# Prompt for directory for cloning

read -p "Enter the directory where you want to clone the repos: " CLONE_DIR

# Create the directory if not present
mkdir -p "$CLONE_DIR"

# Nav to directory
cd "$CLONE_DIR" || exit

# Fetch all repos

repos=$(gh repo list "$GITHUB_USER" --limit 100 --json nameWithOwner --jq '.[].nameWithOwner')

for repo in $repos; do
	git clone "git@github.com:$repo.git"
	echo "$repo has been cloned"
done

echo "All repos have been cloned into $CLONE_DIR"

