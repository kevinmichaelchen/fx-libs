#!/bin/sh

set -e

# Download sd: https://github.com/chmln/sd
brew list -q sd || brew install sd

# Download Gum: https://github.com/charmbracelet/gum
brew list -q gum || brew install gum

# Choose MAJOR, MINOR, or PATCH
TYPE=$(gum choose "patch" "minor" "major")

get_new_tag() {
  # Download SemVer script: https://github.com/ffurrer2/semver
  brew list -q semver || brew install ffurrer2/tap/semver

  LATEST_TAG=$(git tag | sort -V | grep -E '^v\d' | head -n1)
  test -n "$LATEST_TAG" || LATEST_TAG="v0.0.0"
  echo "Most recent tag is: ${LATEST_TAG}"
  # shellcheck disable=SC2016
  # shellcheck disable=SC2039
  LATEST_TAG_WITH_NO_LEADING_V=$(sd '^[v]?(.*)+' '$1' <<< "$LATEST_TAG")

  if [ "$1" = "major" ]; then
    NEW_TAG_WITH_NO_LEADING_V=$(semver next major "${LATEST_TAG_WITH_NO_LEADING_V}")
  elif [ "$1" = "minor" ]; then
    NEW_TAG_WITH_NO_LEADING_V=$(semver next minor "${LATEST_TAG_WITH_NO_LEADING_V}")
  elif [ "$1" = "patch" ]; then
    NEW_TAG_WITH_NO_LEADING_V=$(semver next patch "${LATEST_TAG_WITH_NO_LEADING_V}")
  else
    echo "Illegal SemVer bump type. Should be major|minor|patch"
    exit 1
  fi

  NEW_TAG=v${NEW_TAG_WITH_NO_LEADING_V}
  echo "New tag will be: ${NEW_TAG}"
}

# Confirm with user
gum confirm "Create tags?" && get_new_tag "${TYPE}"

create_tags() {
  PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; | grep -E -v 'tools' | sed 's/^\.\///' | sort)
  COMMIT_HASH=$(git rev-parse --verify HEAD)

  echo "Creating tag for ${NEW_TAG} at commit: ${COMMIT_HASH}"

  # Create tag for root module
  git tag -a "${NEW_TAG}" -m "Version ${NEW_TAG}" "${COMMIT_HASH}"

  # Create tag for submodules
  for dir in $PACKAGE_DIRS; do
    echo "Creating tag for ${dir}/${NEW_TAG} at commit: ${COMMIT_HASH}"
    git tag -a "${dir}/${NEW_TAG}" -m "Version ${dir}/${NEW_TAG}" "${COMMIT_HASH}"
  done

  sleep 2
}

create_tags

# gum spin only works with separate shell scripts
# see https://github.com/charmbracelet/gum/issues/75
#gum spin --spinner dot --title "Buying Bubble Gum..." -- log_tag