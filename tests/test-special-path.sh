#!/usr/bin/env bash
# author: deadc0de6 (https://github.com/deadc0de6)
# Copyright (c) 2024, deadc0de6
#

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

# tmpbasedir
tmpbasedir="/tmp/"
tmpdir=$(mktemp -d "${tmpbasedir}"samedup-specials.XXXX)
clear_on_exit "${tmpdir}"
tmptmp=$(mktemp -d "${tmpbasedir}"samedup-tmp.XXXX)
clear_on_exit "${tmptmp}"

# create a duplicate with space
echo "duplicate1" > "${tmpdir}/a b c.txt"
echo "duplicate1" > "${tmpdir}/de f.txt"

# create a duplicate with special char
echo "duplicate2" > "${tmpdir}/a'bc.txt"
echo "duplicate2" > "${tmpdir}/de'f.txt"

# normal output
"${bin}" "${tmpdir}" | grep '2 duplicates found'

# script output
output="${tmptmp}/script"
"${bin}" --output=script "${tmpdir}" > "${output}"

echo "------------------------"
cat "${output}"
echo "------------------------"

# run generated script
chmod +x "${output}"
"${output}"

# re-run, expect no duplicates
"${bin}" "${tmpdir}" | grep '0 duplicates found'
[ -e "${tmpdir}/de f.txt" ] && echo "not everything removed" && exit 1
[ -e "${tmpdir}/de'f.txt" ] && echo "not everything removed" && exit 1

echo "$(basename "${0}") ok!"
exit 0