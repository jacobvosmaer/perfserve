# perfserve

An all-in-one Go program to capture `perf record` FlameGraphs via a
web browser. Meant to run as root on Linux.

I wanted to be able to get flamegraphs from a Linux machine with minimal effort. With this program I get a single binary that I can quickly copy to a server and run in a `screen` session. 

If you are in a hurry you can skip using `screen` and do:

```
rsync perfserve server.example.com:.
ssh -L localhost:8080:localhost:8080 server.example.com sudo ./perfserve
```

Then visit `http://localhost:8080/?t=5` on your local machine to get a 5-second CPU profile.

## How to build

```
git submodule update --init
GOOS=linux go build
```

This will create an executable `perfserve` that contains embedded
scripts from FlameGraph. If you are compiling on Linux you don't need
the `GOOS=linux` part.

## How to use

Start the server by running `perfserve`. If you want to use a
different HTTP listen address from the default of `localhost:8080`,
use the `-l` flag.

Perfserve depends on `perf` and Perl. Because it runs `perf` it
probably needs to run as `root`.

Perfserve needs a scratch directory where it unpacks its Perl scripts
and where it stores `perf record` data. You can specify this directory
with the `-d` option. If you do not specify a directory, perfserve
will create a temporary directory under the OS tempdir.
