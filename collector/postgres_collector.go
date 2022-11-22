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
	stateDesc        *prometheus.Desc
	roleDesc         *prometheus.Desc
	staticDesc       *prometheus.Desc
	logger           log.Logger
	client           client.PatroniClient
	connectionString string
}

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

func (p *patroniCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.stateDesc
	ch <- p.roleDesc
}

const (
	// Initialize connection constants.
	HOST     = "test-host"
	DATABASE = "test-db"
	USER     = "test-user"
	PASSWORD = "test-creds"
)

func (p *patroniCollector) Collect(ch chan<- prometheus.Metric) {
	// use psql library to run queries
	// pump results to channel

	// connection string
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", HOST, USER, PASSWORD, DATABASE)

	// connection
	db, err := sql.Open("postgres", connectionString)
	checkError(err)

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	/* query

	sql_statement := "SELECT * from ;"
	rows, err := db.Query(sql_statement)
	checkError(err)
	defer rows.Close()

	*/
}
