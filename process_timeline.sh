#!/bin/bash

# Check if jq is installed
if ! command -v jq >/dev/null; then
  echo "Error: jq is not installed. Please install it before running this script."
  exit 1
fi


find . -type f -name "*_*.json"  -exec sh -c '
  for file do
    echo "Processing $file"
    bash process_single_month_timeline.sh "$file" > "${file%.json}_filtered.json"
  done
' exec-sh {} +

