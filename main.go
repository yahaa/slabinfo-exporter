package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yahaa/slabinfo-exporter/slab"
	"golang.org/x/sync/errgroup"
)

const (
	script = `slabtop -o | grep -Ev "^$|OBJS ACTIVE|Minimum / Average / Maximum Object|Active / Total" |awk '{print $1","$2","$3","$4","$5","$6","$7","$8}'`
	port   = 9999
)

var (
	slabInfoCacheSizeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "node",
			Subsystem: "slab_info",
			Name:      "cache_size",
		},
		[]string{"name"},
	)
	slabInfoObjsTotalGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "node",
			Subsystem: "slab_info",
			Name:      "objs_total",
		},
		[]string{"name"},
	)
	slabInfoObjActiveGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "node",
			Subsystem: "slab_info",
			Name:      "objs_active",
		},
		[]string{"name"},
	)

	slabInfoUseGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "node",
			Subsystem: "slab_info",
			Name:      "use",
		},
		[]string{"name"},
	)

	logger = log.With(
		log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout)),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"component", "server",
	)
)

// execWithOutput 运行，且返回运行结果
func execWithOutput(script string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", script)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("run %s err %w", script, err)
	}

	return strings.TrimSpace(string(out)), nil
}

func init() {
	prometheus.MustRegister(slabInfoUseGauge)
	prometheus.MustRegister(slabInfoCacheSizeGauge)
	prometheus.MustRegister(slabInfoObjActiveGauge)
	prometheus.MustRegister(slabInfoObjsTotalGauge)
}

func serve(server *http.Server, port int) func() error {
	return func() error {
		server.Addr = fmt.Sprintf(":%d", port)

		level.Info(logger).Log("msg", fmt.Sprintf("start insecure server for debug, listen on :%d", port))
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("start insecure server err: %w", err)
		}
		return nil
	}
}

func startCollect(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 15)
	for {
		data, err := execWithOutput(script)
		if err != nil {
			logger.Log("msg", fmt.Sprintf("run script err: %v", err))
			return err
		}

		reader := csv.NewReader(strings.NewReader(data))

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Log("msg", fmt.Sprintf("read csv data err: %v", err))
				return err
			}

			s := slab.New(record)

			slabInfoObjsTotalGauge.WithLabelValues(s.Name).Set(s.Objs)
			slabInfoObjActiveGauge.WithLabelValues(s.Name).Set(s.Active)
			slabInfoUseGauge.WithLabelValues(s.Name).Set(s.UseObj)
			slabInfoCacheSizeGauge.WithLabelValues(s.Name).Set(s.CacheSize)
		}

		select {
		case <-ctx.Done():
			level.Info(logger).Log("msg", "stop get slab info...")
			return nil
		case <-ticker.C:
		}
	}
}

func main() {
	var (
		term        = make(chan os.Signal)
		ctx, cancel = context.WithCancel(context.Background())
		r           = mux.NewRouter()
		server      = &http.Server{Handler: r}
	)

	r.Handle("/metrics", promhttp.Handler())

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(serve(server, port))
	wg.Go(func() error { return startCollect(ctx) })

	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
	case <-ctx.Done():
		level.Info(logger).Log("msg", "Stop server, exiting...")
	}

	cancel()

	if err := server.Shutdown(ctx); err != nil {
		level.Error(logger).Log("msg", "Server shutdown error", "err", err)
	}

	if err := wg.Wait(); err != nil {
		level.Error(logger).Log("msg", "Unhandled error received. Exiting...", "err", err)
	}

}
