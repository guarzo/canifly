#!/bin/bash
set -euo pipefail

# Check if jq is installed
if ! command -v jq &>/dev/null; then
    echo "Error: jq is not installed. Please install jq and try again." >&2
    exit 1
fi

# Check if git is installed
if ! command -v git &>/dev/null; then
    echo "Error: git is not installed. Please install git and try again." >&2
    exit 1
fi

version_file="version"
package_file="package.json"
go_file="main.go"

# Ensure version file exists
if [ ! -f "$version_file" ]; then
    echo "Error: $version_file not found. Please create it with an initial version." >&2
    exit 1
fi

current_version=$(cat "$version_file")
echo "Current version: $current_version"

# Increment version
new_version=$(echo "$current_version" | awk -F. -v OFS=. '{ $NF += 1; print }')
echo "New version: $new_version"

# Update the version file
echo "$new_version" > "$version_file"

# Update package.json version
echo "Updating $package_file with new version..."
jq --arg new_version "$new_version" '.version = $new_version' "$package_file" > "$package_file.tmp"
mv "$package_file.tmp" "$package_file"
echo "Updated $package_file."

# Update version in main.go
echo "Updating $go_file with new version..."
sed -i.bak "s/var version = \".*\"/var version = \"$new_version\"/" "$go_file" && rm -f "$go_file.bak"
echo "Updated $go_file."

# Handle optional message
optional_message=${1:-""}
commit_message="chore(release): $new_version"
if [ -n "$optional_message" ]; then
    commit_message="$commit_message - $optional_message"
    echo "Optional message provided: $optional_message"
    # Add all changes if additional message is provided
    git add --all
else
    git add "$version_file" "$package_file" "$go_file"
fi

echo "Committing changes..."
git commit -m "$commit_message"

echo "Tagging new version..."
git tag -a "v$new_version" -m "Release $new_version"

echo "Pushing changes and tags to GitHub..."
git push origin main
git push origin "v$new_version"

echo "Done! New version is $new_version."
