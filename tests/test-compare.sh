#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2023, deadc0de6
#
# depends on
# - fdupes
# - rmlint
# - fclones

set -e

rl="readlink -f"
if ! ${rl} "${0}" >/dev/null 2>&1; then
  rl="realpath"

  if ! hash ${rl}; then
    echo "\"${rl}\" not found !" && exit 1
  fi
fi
cur=$(dirname "$(${rl} "${0}")")
bin="${cur}/../bin/samedup"
hash "${bin}"
# shellcheck source=/dev/null
source "helper.sh"

# ==============
# == the test ==
# ==============

# only output walltime
TIMEFORMAT=%R

# tmpbasedir
tmpbasedir="/tmp/"
tmpbasedir=""

fdupes_args="-r -H -o name"
rmlint_args="-F -o fdupes -r -T df --size 0"
fclones_args="group -A -H --min 0 --hidden"

# $1 path
# $2 other tool and args
# $3 samedup args
# $4 sorter args
compare()
{
  # find duplicates with other tool
  path="${1}"
  echo "\"${2}\" \"${path}\"..."
  time -p ${2} "${path}" 2>/dev/null > "${tmp_other_tool}"
  # sort
  tmp1=$(mktemp "${tmpbasedir}"tmp-samedup-tmp.XXXXX)
  echo "parsing ${tmp_other_tool} to ${tmp1}"
  ./parser.py "${4}" "${tmp_other_tool}" 1> "${tmp1}"
  mv "${tmp1}" "${tmp_other_tool}"

  # find duplicates with samedup
  echo "samedup ${3} \"${path}\"..."
  # shellcheck disable=SC2086
  time -p ${bin} ${3} "${path}" > "${tmp_samedup}"
  # sort
  tmp1=$(mktemp "${tmpbasedir}"tmp-samedup.XXXXX)
  echo "parsing ${tmp_samedup} to ${tmp1}"
  ./parser.py "${4}" "${tmp_samedup}" 1> "${tmp1}"
  mv "${tmp1}" "${tmp_samedup}"

  # diff the results
  echo "diff for \"${path}\" (${cur}/${tmp_other_tool} ${cur}/${tmp_samedup}):"
  if [ -z "${GITHUB_WORKFLOW}" ]; then
    # local
    diff -b -y --suppress-common-lines "${tmp_other_tool}" "${tmp_samedup}"
    #delta --paging=never "${tmp_other_tool}" "${tmp_samedup}"
  else
    # github
    cmp --silent "${tmp_other_tool}" "${tmp_samedup}"
  fi

  echo "ok!"
  echo "--------------------------------------------------------------"
}

# $1 other tool
# $2 other tool args
# $3 path
# $4 samedup args
# $5 sorter args
compares()
{
  local path
  local tool
  local other_args
  local samedup_args
  local sorter_args
  path=${3}
  tool=${1}
  other_args=${2}
  samedup_args=${4}
  sorter_args=${5}
  compare "${path}" "${tool} ${other_args}" "${samedup_args}" "${sorter_args}"
  #compare "/usr" "${tool} ${other_args}" "${samedup_args}" "${sorter_args}"
  compare "/bin" "${tool} ${other_args}" "${samedup_args}" "${sorter_args}"
  compare "/opt" "${tool} ${other_args}" "${samedup_args}" "${sorter_args}"
}

# ensure fdupes exist
hash fdupes
# ensure rmlint exist
hash rmlint
# ensure fclones exist
hash fclones

tmp_other_tool=$(mktemp "${tmpbasedir}"tmp-samedup-other-tool.XXXX)
echo "tmp for other tool: ${cur}/${tmp_other_tool}"
clear_on_exit "${tmp_other_tool}"
tmp_samedup=$(mktemp "${tmpbasedir}"tmp-samedupe.XXXX)
echo "tmp for other tool: ${cur}/${tmp_samedup}"
clear_on_exit "${tmp_samedup}"

path=$(realpath "${cur}"/../.git)

# compare with fdupes
samedup_args="--output=nlines --quiet --sort=name"
compares "fdupes" "${fdupes_args}" "${path}" "${samedup_args}" "--mode=empty-sep"

# compare with rmlint
samedup_args="--output=nlines --quiet --sort=name"
compares "rmlint" "${rmlint_args}" "${path}" "${samedup_args}" "--mode=empty-sep"

# compare with fclones
samedup_args="--output=stairs --quiet --sort=name"
compares "fclones" "${fclones_args}" "${path}" "${samedup_args}" "--mode=stairs"

exit 0