package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/go-chi/chi"
	"github.com/gobuffalo/packr"
	"github.com/markbates/refresh/refresh/web"
	"golang.org/x/net/http2"
)

var (
	fs       = flag.NewFlagSet("app", flag.ExitOnError)
	certFile = fs.String("cert.file", "certs/localhost.cert", "SSL certificate")
	keyFile  = fs.String("key.file", "certs/localhost.key", "Private key")
	httpAddr = fs.String("http.addr", "127.0.0.1:8080", "HTTP server address")
)

func main() {
	fs.Usage = UsageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])
	r := NewRouter()
	var srv http.Server
	srv.Addr = *httpAddr
	srv.Handler = web.ErrorChecker(r)
	http2.ConfigureServer(&srv, nil)
	log.Printf("listening on %s", *httpAddr)
	log.Fatalf("exit with %v", srv.ListenAndServeTLS(*certFile, *keyFile))
}

// NewRouter returns new chi router.
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	box := packr.NewBox("./public")
	root := http.FileServer(box)
	ServeStatic(r, "/", root)
	return r
}

// ServeStatic conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func ServeStatic(r chi.Router, path string, root http.Handler) {
	if strings.ContainsAny(path, "{}*") {
		log.Fatalf("file server does not permit URL parameters")
	}

	fs := http.StripPrefix(path, root)

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// UsageFor is used for print usage.
func UsageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s\t%s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
