#!/bin/sh
set -e

freq=99
sleep=30

perf record -o perf.data -a -g -F $freq -e cpu-clock -- sleep $sleep
perf script -i perf.data | ./stackcollapse-perf.pl | ./flamegraph.pl --title "$1,freq=$freq,sleep=$sleep" --hash
