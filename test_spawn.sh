#!/bin/sh

# Ensure a file path argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <file_path>"
  exit 1
fi

# Get the file path from the first argument
FILE_PATH="$1"

# Get only the current folder name
CURRENT_FOLDER=$(basename "$PWD")

# Write the value of VAR1 and the current folder name to the file
{
  echo "V=${VAR1}"
  echo "D=${CURRENT_FOLDER}"
} > "$FILE_PATH"
