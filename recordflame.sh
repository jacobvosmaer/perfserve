#!/bin/sh
set -e

perf record -o perf.data -a -g -F 99 -- sleep 5
perf script -i perf.data | ./stackcollapse-perf.pl | ./flamegraph.pl --title "$1" --hash
