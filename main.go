package main

import "flag"

var listen = flag.String("l", "localhost:8080", "server listen address")

func main() { flag.Parse() }
