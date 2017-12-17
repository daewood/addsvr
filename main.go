package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"addsvc/strsvc"

	"github.com/go-kit/kit/log"
)

const (
	defaultPort              = "8080"
	defaultRoutingServiceURL = "http://localhost:7878"
)

// curl -XPOST -d'{"s":"hello, world"}' localhost:8000/strsvc/v1/uppercase
// curl -XPOST -d'{"s":"hello, world"}' localhost:8000/strsvc/v1/count
func main() {
	var (
		//		addr = envString("PORT", defaultPort)
		//		rsurl = envString("ROUTINGSERVICE_URL", defaultRoutingServiceURL)
		addr     = "8000"
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		//		routingServiceURL = flag.String("service.routing", rsurl, "routing service URL")

		//		ctx = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	as := strsvc.NewService()
	mux.Handle("/strsvc/v1/", strsvc.MakeHandler(as, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
