[![Go Report Card](https://goreportcard.com/badge/github.com/gopaytech/patroni_exporter)][goreportcard]
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)][license]

[goreportcard]: https://goreportcard.com/report/github.com/thenakedzealot/epg_exporter
[license]: https://opensource.org/licenses/Apache-2.0

# EPG: "Everything PG" Exporter for Prometheus
Scrapes standard Patroni stats and executes custom Postgresql queries deemed essential to monitoring by Aware SRE for active Patroni clusters. Exports metrics via HTTP for Prometheus consumption on port 9933.

## Getting Started

To run it:

```bash
epg_exporter [flags]
```

Help on flags:

```bash
epg_exporter --help
```

For more information check the [source code documentation][gdocs].

[gdocs]: http://godoc.org/github.com/gopaytech/patroni_exporter

## Usage

> Important: Host addresses for both Patroni and Postgres must be supplied for a successful response.

- Specify Patroni API URL using the `--patroni.host` flag.
- Specify Patroni API port using the `--patroni.port` flag.
- Specify Postgres host using the `--postgres.host` flag.
- Specify Postgres user using the `--postgres.user` flag.
- Specify Postgres password using the `--postgres.password` flag.
- Specify Postgres post using the `--postgres.port` flag.
- Specify Postgres database using the `--postgres.database` flag.

```bash
epg_exporter --patroni.host="http://localhost" \
--postgres.host="localhost" --postgres.database="example" \
--postgres.user="superuser" --postgres.password="supersecret" 
```

### Building

```bash
make build
```

### Testing

```bash
make test
```

## License

Apache License 2.0, see [LICENSE](https://github.com/gopaytech/patroni_exporter/blob/master/LICENSE).
