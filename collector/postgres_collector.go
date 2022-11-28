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

// anything need changed here re: connection string? 
type postgresCollector struct {
	stateDesc        *prometheus.Desc
	roleDesc         *prometheus.Desc
	staticDesc       *prometheus.Desc
	logger           log.Logger
	client           client.PatroniClient
	connectionString string
}

// these is creating the descriptions for the metrics? the return of `config.PostgresConnectionString`
func createPostgresCollectorFactory(client client.PatroniClient, config CollectorConfiguration, logger log.Logger) prometheus.Collector {
	stateDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "vacuum_count", "state"),
		"The current vacuum count",
		[]string{"state", "host"},
		nil)
	roleDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "role"),
		"The current PostgreSQL role of Patroni node",
		[]string{"role", "scope"},
		nil)
	staticDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "static"),
		"The collection of static value as reported by Patroni",
		[]string{"version"},
		nil)
	return &patroniCollector{
		stateDesc:        stateDesc,
		roleDesc:         roleDesc,
		staticDesc:       staticDesc,
		logger:           logger,
		client:           client,
		connectionString: config.PostgresConnectionString,
	}
}

	// how do I use?
func (p *patroniCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.stateDesc
	ch <- p.roleDesc
}

// where should go this go, client?
const (
	// Initialize connectionString constants
	HOST     = "test-host"
	DATABASE = "test-db"
	USER     = "test-user"
	PASSWORD = "test-creds"
)


func CheckError(err error) {
	if err != nil {
		panic(err)
}

func (p *patroniCollector) Collect(ch chan<- prometheus.Metric) {
	// use psql library to run queries
	// pump results to channel

	// connection string
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", HOST, USER, PASSWORD, DATABASE)

	// connection attempt
	db, err := sql.Open("postgres", connectionString)
	CheckError(err)

	// success check with ping
	err = db.Ping()
	CheckError(err)
	fmt.Println("Successfully created connection to database")

	// huge query for checkpointlength
	checkpointLength := "SELECT total_checkpoints, seconds_since_start / total_checkpoints / 60 AS minutes_between_checkpoints FROM (SELECT EXTRACT(EPOCH FROM (now() - pg_postmaster_start_time())) AS seconds_since_start,(checkpoints_timed+checkpoints_req) AS total_checkpoints FROM pg_stat_bgwriter) AS sub; SELECT * FROM pg_control_checkpoint();"
	rows, err := db.Query(sql_statement)
	checkError(err)
	defer rows.Close()
	}  //ch <- prometheus.MustNewConstMetric
}


