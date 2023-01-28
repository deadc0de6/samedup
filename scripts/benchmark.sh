#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2023, deadc0de6
#
# usage examples:
#   make && ./scripts/benchmark.sh "${PWD}/bin/samedup --output=stairs" "/usr/bin/fdupes -r" /usr

set -e

# only output walltime
TIMEFORMAT=%R

# $1: bin
# $2: path
run()
{
  bin="${1}"
  path="${2}"
  # push to cache
  for i in {1..3}; do
    echo "caching run ${i} for \"${bin}\""
    ${bin} "${path}" >/dev/null
  done

  # actual tests
  durations=()
  for i in {1..5}; do
    echo ""
    echo "test run ${i} for \"${bin}\""
    real=$( { time ${bin} "${path}" 2>&1 >/dev/null; } 2>&1 )
    echo "duration: ${real}s"
    durations+=( "${real}" )
  done

  echo -e "durations: ${durations[*]}"
  sum=0
  for r in "${durations[@]}"; do
    sum=$(echo "${sum} ${r}" | awk '{print $1 + $2}')
  done
  l=${#durations[@]}
  avg=$(echo "${sum} ${l}" | awk '{print $1 / $2}')
  echo -e "average duration for \"\e[33m${bin}\e[0m\": \e[34m${avg}s\e[0m"
  echo ""
}

usage="${0} <cmd1> <cmd2> <path>"

[ "${1}" = "" ] && echo "${usage}" && exit 1
[ "${2}" = "" ] && echo "${usage}" && exit 1
[ "${3}" = "" ] && echo "${usage}" && exit 1

cmd1="${1}"
cmd2="${2}"
path="${3}"

[ -z "${cmd1}" ] && echo "${cmd1} not found" && exit 1
[ -z "${cmd2}" ] && echo "${cmd2} not found" && exit 1

echo "cmd1: ${cmd1}"
echo "cmd2: ${cmd2}"
echo "path: ${path}"

run "${cmd1}" "${path}"
run "${cmd2}" "${path}"
