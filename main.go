package main

import (
	_ "embed"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var listen = flag.String("l", "localhost:8080", "server listen address")

//go:embed FlameGraph/stackcollapse-perf.pl
var stackcollapse []byte

//go:embed FlameGraph/flamegraph.pl
var flamegraph []byte

//go:embed recordflame.sh
var recordflame []byte

var embeddedFiles map[string][]byte = map[string][]byte{
	"stackcollapse-perf.pl": stackcollapse,
	"flamegraph.pl":         flamegraph,
	"recordflame.sh":        recordflame,
}

func main() {
	flag.Parse()
	if err := run(*listen); err != nil {
		log.Fatal(err)
	}
}

func run(addr string) error {
	dir, err := os.MkdirTemp("", "perfserve")
	if err != nil {
		return err
	}
	os.Chdir(dir)

	for name, data := range embeddedFiles {
		if err := os.WriteFile(name, data, 0755); err != nil {
			return err
		}
	}

	return http.ListenAndServe(addr, &handler{})
}

type handler struct {
	sync.Mutex
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	h.Lock()
	defer h.Unlock()
	rw.Header().Set("Content-Type", "image/svg+xml")
	cmd := exec.Command("./recordflame.sh")
	cmd.Stdout = rw
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		log.Print("recordflame.sh failed: %w", err)
	}

}
