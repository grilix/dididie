package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
)

//go:embed static templates
var files embed.FS

func main() {
	var (
		webAddr = flag.String("web.addr", ":8088", "Web HTTP listen address")
		apiURL  = flag.String("api.url", "http://localhost:8080", "URL where the API dies")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	static, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}
	t, err := template.ParseFS(files, "templates/index.html")
	if err != nil {
		panic(err)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *webAddr)
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/" {
				http.NotFound(w, req)
				return
			}
			w.Header().Add("Content-Type", "text/html")
			t.Execute(w, struct {
				APIUrl string
			}{APIUrl: *apiURL})
		})
		errs <- http.ListenAndServe(*webAddr, nil)
	}()

	logger.Log("exit", <-errs)
}
