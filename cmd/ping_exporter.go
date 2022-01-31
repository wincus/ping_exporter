package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pingexporter "github.com/blainsmith/ping_exporter"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}

func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	metricsAddr := flag.String("metrics.addr", ":9137", "address for ping exporter")
	metricsPath := flag.String("metrics.path", "/metrics", "URL path for surfacing collected metrics")
	interval := flag.Duration("ping.interval", time.Second, "wait time between each ping")

	var pingHosts stringSlice
	flag.Var(&pingHosts, "ping.host", "host to ping, can be repeated (-ping.host=1.1.1.1 -ping.host=google.com ...)")

	flag.Parse()

	if len(pingHosts) <= 0 {
		log.Fatal("no hosts specified to ping")
	}

	var pingers []*ping.Pinger
	for _, host := range pingHosts {
		pinger := ping.New(host)
		pinger.Interval = *interval
		pingers = append(pingers, pinger)

		go func() {
			if err := pinger.Run(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	prometheus.MustRegister(pingexporter.NewCollector(pingers))

	mux := http.NewServeMux()
	mux.Handle(*metricsPath, promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting ping exporter on %q", *metricsAddr)

	if err := http.ListenAndServe(*metricsAddr, mux); err != nil {
		log.Fatalf("cannot start ping exporter: %v", err)
	}
}
