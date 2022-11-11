#!/bin/bash

set -e

PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | egrep -v 'tools' | sed 's/^\.\///' | sort)
COMMIT_HASH=$(git rev-parse --verify HEAD)

echo "Creating tag for ${TAG} at commit: ${COMMIT_HASH}"

# Create tag for root module
git tag -a "${TAG}" -m "Version ${TAG}" ${COMMIT_HASH}

# Create tag for submodules
for dir in $PACKAGE_DIRS; do
	git tag -a "${dir}/${TAG}" -m "Version ${dir}/${TAG}" ${COMMIT_HASH}
done
