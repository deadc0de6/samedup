#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2023, deadc0de6

rl="readlink -f"
if ! ${rl} "${0}" >/dev/null 2>&1; then
  rl="realpath"

  if ! hash ${rl}; then
    echo "\"${rl}\" not found !" && exit 1
  fi
fi
cur=$(dirname "$(${rl} "${0}")")
cwd=$(pwd)

cd "${cur}" || exit 1

for i in "${cur}"/test-*.sh; do
  echo "running ${i}"
  if ! ${i}; then
    echo "test \"${i}\" failed"
    cd "${cwd}" || exit 1
    exit 1
  fi
done

cd "${cwd}" || exit
echo "tests OK!"