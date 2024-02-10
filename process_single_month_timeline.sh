#!/bin/bash

# Check if jq is installed
if ! command -v jq >/dev/null; then
  echo "Error: jq is not installed. Please install it before running this script."
  exit 1
fi

# Get the input file path as an argument
if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <json_file>"
  exit 1
fi

input_file="$1"

# Corrected jq command with proper object construction syntax
output=$(jq '.timelineObjects |= map(
    select(has("placeVisit")
) | .placeVisit.location)' "$input_file")

# Check if jq exited with an error
if [[ $? -ne 0 ]]; then
  echo "Error: jq command failed."
  exit 1
fi

# Print the filtered and deduplicated JSON
echo "$output"

