#!/usr/bin/python3
"""
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6

sort fdupes outputs
"""

import sys
import os
import argparse


UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB']
MODE_STAIRS = 'stairs'
MODE_EMPTY_LINE = 'empty-sep'


class Block:
    """block of duplicates"""
    # pylint: disable=too-few-public-methods

    def __init__(self, entries):
        self.block = entries

    def __lt__(self, other):
        return self.block[0] < other.block[0]


# https://stackoverflow.com/questions/1094841/get-human-readable-version-of-file-size
def human_size(size, units):
    """human size"""
    if size < 1024:
        return str(size) + units[0]
    return human_size(size >> 10, units[1:])


def parse_stairs(file):
    """parse stairs output"""
    dups = []
    acc = []
    cnt = 0
    in_cicd = 'GITHUB_WORKFLOW' in os.environ
    for line in file:
        cnt += 1
        if not in_cicd:
            sys.stderr.write(f'line {cnt}\r')
        if line == '':
            continue
        if line.startswith('#'):
            continue
        if line.startswith(' ') or line.startswith('\t'):
            # append
            line = line.strip()
            path = os.path.realpath(line)
            acc.append(path)
        else:
            # separator
            if len(acc) > 0:
                acc.sort()
                dups.append(Block(acc))
                acc = []
    if len(acc) > 0:
        acc.sort()
        dups.append(Block(acc))
    sys.stderr.write(f'found {len(dups)} block(s) on {cnt} line(s)\n')
    return dups


def parse_empty_sep(file):
    """parse empty line separated output"""
    dups = []
    acc = []
    cnt = 0
    for line in file:
        cnt += 1
        sys.stderr.write(f'line {cnt}\r')
        line = line.rstrip()
        if line == '':
            # separator
            if len(acc) > 0:
                acc.sort()
                dups.append(Block(acc))
                acc = []
        else:
            # append
            path = os.path.realpath(line)
            acc.append(path)
    if len(acc) > 0:
        acc.sort()
        dups.append(Block(acc))
    sys.stderr.write(f'found {len(dups)} block(s) on {cnt} line(s)\n')
    return dups


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('path')
    parser.add_argument('--mode', choices=[MODE_EMPTY_LINE, MODE_STAIRS])
    args = parser.parse_args()

    file_size = human_size(os.path.getsize(args.path), UNITS)
    sys.stderr.write(
        f'parsing and sorting \"{args.path}\" '
        f'(size: {file_size})\n'
        )

    blocks = []
    with open(args.path, 'r', encoding="utf-8") as filed:
        if args.mode == MODE_EMPTY_LINE:
            blocks = parse_empty_sep(filed)
        elif args.mode == MODE_STAIRS:
            blocks = parse_stairs(filed)

    blocks.sort()
    for block in blocks:
        print('\n'.join(block.block))
        print('')
