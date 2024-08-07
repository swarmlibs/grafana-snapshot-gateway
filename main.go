package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gofrs/uuid"
	"github.com/prometheus-operator/prometheus-operator/pkg/versionutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/swarmlibs/grafana-snapshot-gateway/grafana"
	"github.com/swarmlibs/grafana-snapshot-gateway/grafana/types"
	"github.com/swarmlibs/grafana-snapshot-gateway/internal/metrics"
	"github.com/swarmlibs/grafana-snapshot-gateway/internal/middlewares"
)

func main() {
	app := kingpin.New("grafana-snapshot-gateway", "")

	// Server options
	listenAddr := app.Flag("listen-addr", "The address to listen on for HTTP requests.").Default(":3003").String()

	// Grafana snapshot server
	grafanaUrl := app.Flag("grafana-url", "Grafana URL").Required().String()
	grafanaBasicAuth := app.Flag("grafana-basic-auth", "Grafana credentials").String()

	// Snapshot deletion options
	checkSnapshotBeforeDelete := app.Flag("check-snapshot-before-delete", "Check if snapshot exists before delete").Default("false").Bool()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stdout)
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	versionutil.RegisterIntoKingpinFlags(app)

	if versionutil.ShouldPrintVersion() {
		versionutil.Print(os.Stdout, "grafana-snapshot-gateway")
		os.Exit(0)
	}

	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(2)
	}

	ps := prometheus.NewRegistry()
	ps.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	mc := metrics.New("gf_snapshot_gateway", ps)

	level.Info(logger).Log("msg", "Starting node-metadata-agent", "version", version.Info())

	level.Info(logger).Log("msg", fmt.Sprintf("Listening on %s", *listenAddr))
	level.Info(logger).Log("msg", fmt.Sprintf("Grafana URL: %s", *grafanaUrl))

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.StructuredLogger(&logger))
	r.Use(middlewares.MeasureRequestDuration(mc))
	r.SetTrustedProxies(nil)

	gf := grafana.NewGrafanaClient(*grafanaUrl, "", "")
	gf.SetLogger(logger)

	if *grafanaBasicAuth != "" {
		creds := strings.Split(*grafanaBasicAuth, ":")
		if len(creds) != 2 {
			level.Error(logger).Log("msg", "Invalid credentials")
			os.Exit(2)
		}
		gf.SetBasicAuth(creds[0], creds[1])
		level.Info(logger).Log("msg", "Grafana basic auth enabled", "uname", creds[0])
	}

	if *checkSnapshotBeforeDelete {
		level.Info(logger).Log("msg", "Check snapshot before delete enabled")
	}

	level.Info(logger).Log("build_context", version.BuildContext())

	// ANY /
	// Redirect to Grafana URL
	r.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, *grafanaUrl)
	})

	// GET /health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// GET /metrics
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(ps, promhttp.HandlerOpts{})))

	// Create new snapshot
	// POST /api/snapshots
	r.POST("/api/snapshots", func(c *gin.Context) {
		var err error
		var snapshot types.GrafanaDashboardSnapshot

		if err := c.ShouldBindJSON(&snapshot); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Create a new overrideUid for the folder, dashboard and snapshot
		originalUid := snapshot.Dashboard.Get("uid")
		overrideUid := uuid.Must(uuid.NewV4()).String()

		snapshot.SetKey(overrideUid)
		level.Info(logger).Log("msg", "Create new snapshot", "uid", originalUid, "uid_overrided", overrideUid)

		// Create a new folder
		_, err = gf.CreateFolder(overrideUid, overrideUid)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Create a new snapshot
		snapshotCreationResponse, err := gf.CreateSnapshot(snapshot.Key, snapshot)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Delete the folder
		gf.DeleteFolder(snapshot.Key)

		// Customize the snapshot response
		snapshotResponse := types.GrafanaDashboardSnapshotCreateResponse{}
		grafana.UnmarshalResponseBody(snapshotCreationResponse.Body, &snapshotResponse)

		// If possible detect "Host" header and override deleteUrl
		// The reason is that if Grafana is able to talk to the gateway, it should be use the gateway to delete the snapshot
		// instead of the Grafana instance
		if host := c.Request.Host; host != "" {
			oldDeleteUrl := snapshotResponse.DeleteUrl
			snapshotResponse.SetDeleteUrlHost(host)
			level.Info(logger).Log("msg", "Override snapshot delete url", "delete_url", oldDeleteUrl, "delete_url_overrided", snapshotResponse.DeleteUrl)
		}

		// Return the snapshot response
		level.Info(logger).Log("msg", "Snapshot created successfully", "uid", originalUid, "uid_overrided", overrideUid)
		c.JSON(snapshotCreationResponse.StatusCode, snapshotResponse)
	})

	// Delete Snapshot by Key
	// GET /api/snapshots-delete/:key
	r.GET("/api/snapshots-delete/:key", func(c *gin.Context) {
		key := c.Param("key")
		level.Info(logger).Log("msg", "Delete snapshot", "key", key)
		res, err := gf.DeleteSnapshot(key)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var response types.GrafanaDashboardSnapshotDeleteResponse
		if !*checkSnapshotBeforeDelete {
			response.Message = "Snapshot deleted successfully"
			c.JSON(http.StatusOK, response)
		} else {
			grafana.UnmarshalResponseBody(res.Body, &response)
			c.JSON(res.StatusCode, response)
		}

		level.Info(logger).Log("msg", response.Message, "key", key)
	})

	// listen and serve, default 0.0.0.0:3003 (for windows "localhost:3003")
	if err := r.Run(*listenAddr); err != nil {
		level.Error(logger).Log("msg", "Failed to start server", "err", err)
		os.Exit(1)
	}
}
