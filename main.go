package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
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
	if err := os.Chmod(dir, 0700); err != nil {
		return err
	}
	if err := os.Chdir(dir); err != nil {
		return err
	}

	for name, data := range embeddedFiles {
		if err := os.WriteFile(name, data, 0755); err != nil {
			return err
		}
	}
	h := &handler{}
	if h.hostname, err = os.Hostname(); err != nil {
		return err
	}

	return http.ListenAndServe(addr, h)
}

type handler struct {
	hostname string
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		return
	}
	freq := getPositiveInt(r, "f", 99)
	sleep := getPositiveInt(r, "t", 30)

	now := time.Now().UTC()
	y, m, d := now.Date()
	title := fmt.Sprintf("%s-%d-%02d-%02dT%02d_%02d.%03dZ", h.hostname, y, m, d, now.Hour(), now.Minute(), now.UnixMilli()%1000)

	hdr := rw.Header()
	hdr.Set("Content-Type", "image/svg+xml")
	hdr.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.svg"`, title))
	cmd := exec.Command("./recordflame.sh", title, strconv.Itoa(freq), strconv.Itoa(sleep))
	cmd.Stdout = rw
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		log.Print("recordflame.sh failed: %w", err)
	}
}

func getPositiveInt(r *http.Request, field string, defaultValue int) int {
	val, _ := strconv.ParseInt(r.Form.Get(field), 10, 32)
	if val <= 0 {
		val = int64(defaultValue)
	}
	return int(val)
}
