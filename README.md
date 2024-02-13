# perfserve

An all-in-one Go program to capture `perf record` FlameGraphs via a web browser. Meant to run as root on Linux.

I'm sure there exist proper tools for this but I felt like building this myself.

## How to build

```
git submodule update --init
GOOS=linux go build
```

This will create an executable `perfserve` that contains embedded scripts from FlameGraph. If you are compiling on Linux you don't need the `GOOS=linux` part.

## How to use

Start the server by running `perfserve`. If you want to use a different HTTP listen address from the default of `localhost:8080`, use the `-l` flag.

The server depends on `perf` and Perl. Because it runs `perf` it probably needs to run as `root`.

When `perfserve` starts up it will create a temporary directory, unpack its Perl scripts in there, and it will use that directory to store `perf record` data.

