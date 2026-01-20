# Environment Variables

This document describes the environment variables used by `buang`.

| Name                            | Usage                                                       | Description                                                                                                   |
| ------------------------------- | ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| [`API_KEY`]                     | defaults to `supersecureapikey`                             | Buang API key                                                                                                 |
| [`DBOS_ADMIN_SERVER_PORT`]      | optional                                                    | DBOS admin server port. Optional to use DBOS. Specify a port to enable the admin server, DBOS default is 3001 |
| [`DBOS_CONDUCTOR_API_KEY`]      | optional                                                    | DBOS conductor API key. Optional to use DBOS                                                                  |
| [`DBOS_CONDUCTOR_URL`]          | optional                                                    | DBOS conductor url. Optional to use DBOS. You can use either the DBOS console or self host it                 |
| [`DB_CONNECTION_STRING`]        | defaults to `'file:test.db?_foreign_keys=true&mode=memory'` | Database connection string                                                                                    |
| [`DB_TYPE`]                     | defaults to `sqlite`                                        | Database type                                                                                                 |
| [`DURABLE_EXECUTOR`]            | defaults to `temporal`                                      | Durable executor                                                                                              |
| [`ENABLE_TELEMETRY`]            | defaults to `false`                                         | Enable telemetry                                                                                              |
| [`GO_ENV`]                      | defaults to `development`                                   | Golang environment                                                                                            |
| [`LOG_LEVEL`]                   | defaults to `INFO`                                          | Log level                                                                                                     |
| [`OTEL_EXPORTER_OTLP_ENDPOINT`] | defaults to `''`                                            | OpenTelemetry exporter endpoint                                                                               |
| [`OTEL_SERVICE_NAME`]           | defaults to `skulpture-buang`                               | OpenTelemetry service name. This is also the DBOS application name which is required when using DBOS          |
| [`TEMPORAL_ADDRESS`]            | optional                                                    | Temporal address. You can either use Temporal Cloud or self host it                                           |
| [`TEMPORAL_API_KEY`]            | optional                                                    | Temporal API key                                                                                              |
| [`TEMPORAL_NAMESPACE`]          | optional                                                    | Temporal namespace                                                                                            |

> [!TIP]
> If an environment variable is set to an empty value, `buang` behaves as if
> that variable is left undefined.

## `API_KEY`

> Buang API key

The `API_KEY` variable **MAY** be left undefined, in which case a default value
is used.

⚠️ This variable is **sensitive**; its value may contain private information.

## `DBOS_ADMIN_SERVER_PORT`

> DBOS admin server port. Optional to use DBOS. Specify a port to enable the admin server, DBOS default is 3001

The `DBOS_ADMIN_SERVER_PORT` variable **MAY** be left undefined. Otherwise, the
value **MUST** be a whole number.

```bash
export DBOS_ADMIN_SERVER_PORT=-922337203685477888  # (non-normative)
export DBOS_ADMIN_SERVER_PORT=+1844674407370954752 # (non-normative)
```

<details>
<summary>Signed integer syntax</summary>

Signed integers can only be specified using decimal notation. A leading positive
sign (`+`) is **OPTIONAL**. A leading negative sign (`-`) is **REQUIRED** in
order to specify a negative value.

Internally, the `DBOS_ADMIN_SERVER_PORT` variable is represented using a signed
64-bit integer type (`int`); any value that overflows this data-type is invalid.

</details>

## `DBOS_CONDUCTOR_API_KEY`

> DBOS conductor API key. Optional to use DBOS

The `DBOS_CONDUCTOR_API_KEY` variable **MAY** be left undefined.

⚠️ This variable is **sensitive**; its value may contain private information.

## `DBOS_CONDUCTOR_URL`

> DBOS conductor url. Optional to use DBOS. You can use either the DBOS console or self host it

The `DBOS_CONDUCTOR_URL` variable **MAY** be left undefined.

```bash
export DBOS_CONDUCTOR_URL=foo # (non-normative)
```

## `DB_CONNECTION_STRING`

> Database connection string

The `DB_CONNECTION_STRING` variable **MAY** be left undefined, in which case a
default value is used.

⚠️ This variable is **sensitive**; its value may contain private information.

## `DB_TYPE`

> Database type

The `DB_TYPE` variable **MAY** be left undefined, in which case the default
value of `sqlite` is used. Otherwise, the value **MUST** be either `postgres` or
`sqlite`.

```bash
export DB_TYPE=postgres
export DB_TYPE=sqlite   # (default)
```

## `DURABLE_EXECUTOR`

> Durable executor

The `DURABLE_EXECUTOR` variable **MAY** be left undefined, in which case the
default value of `temporal` is used. Otherwise, the value **MUST** be either
`temporal` or `dbos`.

```bash
export DURABLE_EXECUTOR=temporal # (default)
export DURABLE_EXECUTOR=dbos
```

## `ENABLE_TELEMETRY`

> Enable telemetry

The `ENABLE_TELEMETRY` variable **MAY** be left undefined, in which case the
default value of `false` is used. Otherwise, the value **MUST** be either `true`
or `false`.

```bash
export ENABLE_TELEMETRY=true
export ENABLE_TELEMETRY=false # (default)
```

## `GO_ENV`

> Golang environment

The `GO_ENV` variable **MAY** be left undefined, in which case the default value
of `development` is used. Otherwise, the value **MUST** be one of the values
shown in the examples below.

```bash
export GO_ENV=production
export GO_ENV=development # (default)
export GO_ENV=test
```

## `LOG_LEVEL`

> Log level

The `LOG_LEVEL` variable **MAY** be left undefined, in which case the default
value of `INFO` is used. Otherwise, the value **MUST** be one of the values
shown in the examples below.

```bash
export LOG_LEVEL=DEBUG
export LOG_LEVEL=ERROR
export LOG_LEVEL=INFO  # (default)
export LOG_LEVEL=WARN
```

## `OTEL_EXPORTER_OTLP_ENDPOINT`

> OpenTelemetry exporter endpoint

The `OTEL_EXPORTER_OTLP_ENDPOINT` variable **MAY** be left undefined, in which
case the default value of `` is used.

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT='' # (default)
```

## `OTEL_SERVICE_NAME`

> OpenTelemetry service name. This is also the DBOS application name which is required when using DBOS

The `OTEL_SERVICE_NAME` variable **MAY** be left undefined, in which case the
default value of `skulpture-buang` is used.

```bash
export OTEL_SERVICE_NAME=skulpture-buang # (default)
```

## `TEMPORAL_ADDRESS`

> Temporal address. You can either use Temporal Cloud or self host it

The `TEMPORAL_ADDRESS` variable **MAY** be left undefined.

```bash
export TEMPORAL_ADDRESS=foo # (non-normative)
```

## `TEMPORAL_API_KEY`

> Temporal API key

The `TEMPORAL_API_KEY` variable **MAY** be left undefined.

⚠️ This variable is **sensitive**; its value may contain private information.

## `TEMPORAL_NAMESPACE`

> Temporal namespace

The `TEMPORAL_NAMESPACE` variable **MAY** be left undefined.

```bash
export TEMPORAL_NAMESPACE=foo # (non-normative)
```

---

> [!NOTE]
> This document only describes environment variables declared using [Ferrite].
> `buang` may consume other undocumented environment variables.

> [!IMPORTANT]
> Some of the example values given in this document are **non-normative**.
> Although these values are syntactically valid, they may not be meaningful to
> `buang`.

<!-- references -->

[`api_key`]: #API_KEY
[`db_connection_string`]: #DB_CONNECTION_STRING
[`db_type`]: #DB_TYPE
[`dbos_admin_server_port`]: #DBOS_ADMIN_SERVER_PORT
[`dbos_conductor_api_key`]: #DBOS_CONDUCTOR_API_KEY
[`dbos_conductor_url`]: #DBOS_CONDUCTOR_URL
[`durable_executor`]: #DURABLE_EXECUTOR
[`enable_telemetry`]: #ENABLE_TELEMETRY
[ferrite]: https://github.com/dogmatiq/ferrite
[`go_env`]: #GO_ENV
[`log_level`]: #LOG_LEVEL
[`otel_exporter_otlp_endpoint`]: #OTEL_EXPORTER_OTLP_ENDPOINT
[`otel_service_name`]: #OTEL_SERVICE_NAME
[`temporal_address`]: #TEMPORAL_ADDRESS
[`temporal_api_key`]: #TEMPORAL_API_KEY
[`temporal_namespace`]: #TEMPORAL_NAMESPACE
