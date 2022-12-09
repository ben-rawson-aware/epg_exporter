package collector

import (
	"fmt"

	"github.com/go-kit/kit/log"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/gopaytech/patroni_exporter/client"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("postgres", createPatroniCollectorFactory)
}

var (
	possiblePatroniState = [...]string{"RUNNING", "STOPPED", "PROMOTED", "UNKNOWN"}
	possiblePatroniRole  = [...]string{"MASTER", "REPLICA", "STANDBY_LEADER"}
)

type postgresCollector struct {
	vacuumDesc       *prometheus.Desc
	activeDesc       *prometheus.Desc
	connectDesc      *prometheus.Desc
	baselineDesc     *prometheus.Desc
	checkpointDesc   *prometheus.Desc
	logger           log.Logger
	client           client.PatroniClient
	config           CollectorConfiguration
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
	return &patroniCollector{
		vacuumDesc:       vacuumDesc,
		connectDesc:      connectDesc,
		baselineDesc:     baselineDesc,
		checkpointDesc:   checkpointDesc,
		logger:           logger,
		client:           client,
		config: 		  config,
	}
}

func (p *patroniCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.vacuumDesc
	ch <- p.connectDesc
	ch <- p.baselineDesc
	ch <- p.checkpointDesc
}

func CheckError(err error) {
	if err != nil {
		panic(err)
}

func (p *patroniCollector) Collect(ch chan<- prometheus.Metric) {

	// connection string
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", p.config.HOST, p.config.USER, p.config.PASSWORD, p.config.DATABASE)

	// connection attempt
	db, err := sql.Open("postgres", connectionString)
	CheckError(err)

	// success check with ping
	err = db.Ping()
	CheckError(err)
	fmt.Println("Successfully created connection to database")

	// vacuum activity
	vacuum_act := "SELECT count(pid) FROM pg_stat_activity WHERE query LIKE 'autovacuum: %';"
	rows, err := db.Query(vacuum_act)
	checkError(err)
	ch <- prometheus.MustNewConstMetric(p.vacuumDesc, prometheus.CounterValue, prometheus.CounterValue, stateValue, possibleRole, patroniResponse.Patroni.Scope)
	rows.Close()

	// active connection count (need to redo query for groupby state)
	connect_count := "SELECT count(*) FROM pg_stat_activity;"
	rows, err := db.Query(connect_count)
	checkError(err)
	ch <- prometheus.MustNewConstMetric(p.connectDesc, prometheus.CounterValue, prometheus.CounterValue, stateValue, possibleRole, patroniResponse.Patroni.Scope)
	rows.Close()

	// baseline transaction rates
	baseline_sum := "SELECT sum(xact_commit+xact_rollback) FROM pg_stat_database;"
	rows, err := db.Query(baseline_sum)
	checkError(err)
	ch <- prometheus.MustNewConstMetric(p.baselineDesc, prometheus.CounterValue, prometheus.CounterValue, stateValue, possibleRole, patroniResponse.Patroni.Scope)
	rows.Close()

	// huge query for checkpointlength (checkpoints_last60min = total_checkpoints)avg this
	checkpoint_length := "SELECT avg(checkpoints_timed+checkpoints_req) AS total_checkpoints FROM pg_stat_bgwriter WHERE time BETWEEN now() - '1 hour' AND now();"
	rows, err := db.Query(checkpoint_length)
	checkError(err)
	ch <- prometheus.MustNewConstMetric(p.checkpointDesc, prometheus.CounterValue, prometheus.CounterValue, stateValue, possibleRole, patroniResponse.Patroni.Scope)
	rows.Close()

}
