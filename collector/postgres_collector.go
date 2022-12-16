package collector

import (
	"os"

	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/gopaytech/patroni_exporter/client"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("postgres", createPostgresCollectorFactory)
}

type postgresCollector struct {
	vacuumDesc     *prometheus.Desc
	connectDesc    *prometheus.Desc
	baselineDesc   *prometheus.Desc
	checkpointDesc *prometheus.Desc
	logger         log.Logger
	client         client.PatroniClient
	config         CollectorConfiguration
}

func createPostgresCollectorFactory(client client.PatroniClient, config CollectorConfiguration, logger log.Logger) prometheus.Collector {
	vacuumDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "vacuum"),
		"The current vacuum activity",
		[]string{"cluster_name", "machine_name", "role"},
		nil)
	connectDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "connect"),
		"The current connection count",
		[]string{"cluster_name", "machine_name", "role"},
		nil)
	baselineDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "baseline"),
		"The current baseline transaction rate",
		[]string{"cluster_name", "machine_name", "role"},
		nil)
	checkpointDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "checkpoint"),
		"Checkpoint average for last hour",
		[]string{"cluster_name", "machine_name", "role"},
		nil)
	return &postgresCollector{
		vacuumDesc:     vacuumDesc,
		connectDesc:    connectDesc,
		baselineDesc:   baselineDesc,
		checkpointDesc: checkpointDesc,
		logger:         logger,
		client:         client,
		config:         config,
	}
}

func (p *postgresCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.vacuumDesc
	ch <- p.connectDesc
	ch <- p.baselineDesc
	ch <- p.checkpointDesc
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func (p *postgresCollector) Collect(ch chan<- prometheus.Metric) {

	patroniResponse, err := p.client.GetMetrics()
	if err != nil {
		level.Error(p.logger).Log("msg", "Unable to get metrics from Patroni", "err", fmt.Sprintf("errornya: %v", err))
		return
	}

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// connection string
	var connectionString string = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", p.config.HOST, p.config.PORT, p.config.USER, p.config.PASSWORD, p.config.DATABASE)

	// connection attempt
	db, err := sql.Open("postgres", connectionString)
	CheckError(err)

	// success check with ping
	err = db.Ping()
	CheckError(err)
	fmt.Println("Successfully created connection to database")

	// vacuum activity
	var count int
	vacuum_act := "SELECT count(pid) FROM pg_stat_activity WHERE query LIKE 'autovacuum: %';"
	row := db.QueryRow(vacuum_act)
	err = row.Scan(&count)
	CheckError(err)
	ch <- prometheus.MustNewConstMetric(p.vacuumDesc, prometheus.CounterValue, float64(count), p.config.CLUSTER, name, patroniResponse.Role)

	// active connection count grouped by state
	connect_count := "SELECT count(*) FROM pg_stat_activity GROUP BY state;"
	row = db.QueryRow(connect_count)
	err = row.Scan(&count)
	CheckError(err)
	ch <- prometheus.MustNewConstMetric(p.connectDesc, prometheus.CounterValue, float64(count), p.config.CLUSTER, name, patroniResponse.Role)

	// baseline transaction rates
	var sum float64
	baseline_sum := "SELECT sum(xact_commit+xact_rollback) FROM pg_stat_database;"
	row = db.QueryRow(baseline_sum)
	err = row.Scan(&sum)
	CheckError(err)
	ch <- prometheus.MustNewConstMetric(p.baselineDesc, prometheus.CounterValue, sum, p.config.CLUSTER, name, patroniResponse.Role)

	// huge query for checkpointlength (checkpoints_last60min = total_checkpoints)avg this
	checkpoint_length := "SELECT avg(seconds_since_start / total_checkpoints / 60) AS checkpoint_length FROM (SELECT EXTRACT(EPOCH FROM (now() - pg_postmaster_start_time())) AS seconds_since_start, (checkpoints_timed+checkpoints_req) AS total_checkpoints FROM pg_stat_bgwriter) AS sub;"
	row = db.QueryRow(checkpoint_length)
	err = row.Scan(&sum)
	CheckError(err)
	ch <- prometheus.MustNewConstMetric(p.checkpointDesc, prometheus.CounterValue, sum, p.config.CLUSTER, name, patroniResponse.Role)
}
