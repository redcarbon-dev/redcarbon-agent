#!/bin/bash

user_id=`id -u`

# we want to snapshot the environment of the config user
if [ $user_id -eq 0 -a -z "$AGENT_ALLOW_RUNASROOT" ]; then
    echo "Must not run with sudo"
    exit 1
fi

# Change directory to the script root directory
# https://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
cd "$DIR"

shopt -s nocasematch
./bin/redcarbon profile "$@"
