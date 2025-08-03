#!/bin/bash
# Simple release script for manual version management
set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if version file exists
if [ ! -f "version" ]; then
    echo -e "${RED}Error: version file not found${NC}"
    exit 1
fi

# Read current version
version=$(cat version)
echo -e "${GREEN}Current version: $version${NC}"

# Check if package.json version matches
package_version=$(grep -o '"version": "[^"]*"' package.json | cut -d'"' -f4)
if [ "$version" != "$package_version" ]; then
    echo -e "${YELLOW}Warning: package.json version ($package_version) doesn't match version file ($version)${NC}"
    echo "Please update package.json manually to match"
    exit 1
fi

# Optional commit message
message=${1:-""}

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}You have uncommitted changes. Please commit or stash them first.${NC}"
    exit 1
fi

# Create git tag
tag="v$version"
echo -e "${GREEN}Creating tag: $tag${NC}"

if [ -n "$message" ]; then
    git tag -a "$tag" -m "Release $version: $message"
else
    git tag -a "$tag" -m "Release $version"
fi

echo -e "${GREEN}Tag created successfully!${NC}"
echo ""
echo "To push this release:"
echo "  git push origin main"
echo "  git push origin $tag"
echo ""
echo "To update the version for next release:"
echo "  1. Edit the 'version' file"
echo "  2. Update 'version' in package.json to match"
echo "  3. Commit both changes"
echo "  4. Run this script again"