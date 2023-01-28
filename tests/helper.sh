#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2023, deadc0de6

declare -a to_be_cleared

# add a file/directory to be cleared
# on exit
#
# $1: file path to clear
clear_on_exit()
{
  local len="${#to_be_cleared[*]}"
  to_be_cleared["${len}"]="$1"
  if [ "${len}" = "0" ]; then
    # set trap
    trap on_exit EXIT
  fi
}

# clear files
on_exit()
{
  for i in "${to_be_cleared[@]}"; do
    rm -rf "${i}"
  done
}

