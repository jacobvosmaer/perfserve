#!/bin/sh
set -e
perf record -o perf.data -a -g -F 99 -- sleep 30
perf script -i perf.data | ./stackcollapse-perf.pl | ./flamegraph.pl --title "$(hostname)" --hash
