package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gofrs/uuid"
	"github.com/prometheus-operator/prometheus-operator/pkg/versionutil"
	"github.com/prometheus/common/version"
	"github.com/swarmlibs/grafana-snapshot-gateway/grafana"
)

type GrafanaSnapshotBody struct {
	Dashboard map[string]interface{} `json:"dashboard"`
	Name      string                 `json:"name,omitempty"`
	Expires   int                    `json:"expires,omitempty"`
	External  bool                   `json:"external,omitempty"`
	Key       string                 `json:"key,omitempty"`
	DeleteKey string                 `json:"deleteKey,omitempty"`
}

func main() {
	app := kingpin.New("grafana-snapshot-gateway", "")

	listenAddr := app.Flag("listen-addr", "The address to listen on for HTTP requests.").Default(":3003").String()
	grafanaUrl := app.Flag("grafana-url", "Grafana URL").Required().String()
	grafanaBasicAuth := app.Flag("grafana-basic-auth", "Grafana credentials").String()

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

	grafanaClient := grafana.NewGrafanaClient(*grafanaUrl, "", "")
	if *grafanaBasicAuth != "" {
		creds := strings.Split(*grafanaBasicAuth, ":")
		if len(creds) != 2 {
			level.Error(logger).Log("msg", "Invalid credentials")
			os.Exit(1)
		}
		grafanaClient.SetBasicAuth(creds[0], creds[1])
	}

	// POST /api/snapshots
	r.POST("/api/snapshots", func(c *gin.Context) {
		var body GrafanaSnapshotBody
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if body.Dashboard == nil {
			c.JSON(400, gin.H{"error": "dashboard is required"})
			return
		}

		uid := uuid.Must(uuid.NewV4()).String()
		grafanaClient.CreateFolder(uid, uid)

		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// listen and serve, default 0.0.0.0:3003 (for windows "localhost:3003")
	r.Run(*listenAddr)
}
