/////// This should be something else

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
)

func main() {
	var (
		webAddr = flag.String("web.addr", ":8088", "Web HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	fs := http.FileServer(http.Dir("./cmd/web/public"))

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *webAddr)
		http.Handle("/", fs)
		errs <- http.ListenAndServe(*webAddr, nil)
	}()

	logger.Log("exit", <-errs)
}
