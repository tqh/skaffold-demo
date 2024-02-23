package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	counter   = make(map[string]int64)
	countLock sync.Mutex
)

func count(name string) int64 {
	countLock.Lock()
	defer countLock.Unlock()
	counter[name] += 1
	return counter[name]
}

func nameFromPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	return html.EscapeString(path)
}

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		name := nameFromPath(r.URL.Path)
		_, _ = fmt.Fprintf(w, "Hi, %q: %d", name, pgcount(name))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
