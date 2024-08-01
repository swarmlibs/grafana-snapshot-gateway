package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus-operator/prometheus-operator/pkg/versionutil"
	"github.com/prometheus/common/version"
)

func main() {
	app := kingpin.New("grafana-snapshot-gateway", "")

	listenAddr := app.Flag("listen-addr", "The address to listen on for HTTP requests.").Default(":3003").String()
	grafanaUrl := app.Flag("grafana-url", "Grafana URL").Required().String()
	// grafanaCredentials := app.Flag("grafana-credentials", "Grafana credentials").String()

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

	level.Info(logger).Log("msg", "Starting node-metadata-agent", "version", version.Info())
	level.Info(logger).Log("msg", fmt.Sprintf("Listening on %s", *listenAddr))
	level.Info(logger).Log("msg", fmt.Sprintf("Grafana URL: %s", *grafanaUrl))
	level.Info(logger).Log("build_context", version.BuildContext())

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.Default()
	r.SetTrustedProxies(nil)

	// GET /api/snapshots
	r.GET("/api/snapshots", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(*listenAddr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
