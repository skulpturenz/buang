# Environment Variables

This document describes the environment variables used by `buang`.

| Name                                            | Usage                                                                    | Description                                                                                                                    |
| ----------------------------------------------- | ------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------ |
| [`BUANG_API_KEY`]                               | defaults to `supersecureapikey`                                          | Buang API key                                                                                                                  |
| [`BUANG_DBOS_ADMIN_SERVER_PORT`]                | optional                                                                 | DBOS admin server port. Optional to use DBOS. Specify a port to enable the admin server, DBOS default is 3001                  |
| [`BUANG_DBOS_CONDUCTOR_API_KEY`]                | optional                                                                 | DBOS conductor API key. Optional to use DBOS                                                                                   |
| [`BUANG_DBOS_CONDUCTOR_URL`]                    | optional                                                                 | DBOS conductor url. Optional to use DBOS. You can use either the DBOS console or self host it                                  |
| [`BUANG_DB_CONNECTION_STRING`]                  | defaults to `'file:test.db?_foreign_keys=true&mode=memory&cache=shared'` | Database connection string                                                                                                     |
| [`BUANG_DB_TYPE`]                               | defaults to `sqlite`                                                     | Database type                                                                                                                  |
| [`BUANG_DURABLE_EXECUTOR`]                      | defaults to `temporal`                                                   | Durable executor                                                                                                               |
| [`BUANG_ENABLE_EXPERIMENTAL_BOOTSTRAP`]         | defaults to `false`                                                      | Enable experimental bootstrap. Bootstrapping allows Buang to deploy itself and autoupdates on Saturdays at midnight every week |
| [`BUANG_ENABLE_TELEMETRY`]                      | defaults to `false`                                                      | Enable telemetry                                                                                                               |
| [`BUANG_EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE`] | defaults to `'0 0 * * 6'`                                                | Cron expression for when auto updates should run. Defaults to every Saturday at midnight                                       |
| [`BUANG_HOUSEKEEPING_PRUNE_DEPLOYMENTS`]        | defaults to `'0 0 */2 * *'`                                              | Cron expression for stale deployments should be pruned. Defaults to every 2 days                                               |
| [`BUANG_LOG_LEVEL`]                             | defaults to `INFO`                                                       | Log level                                                                                                                      |
| [`BUANG_OTEL_EXPORTER_OTLP_ENDPOINT`]           | defaults to `''`                                                         | OpenTelemetry exporter endpoint                                                                                                |
| [`BUANG_OTEL_SERVICE_NAME`]                     | defaults to `skulpture-buang`                                            | OpenTelemetry service name. This is also the DBOS application name which is required when using DBOS                           |
| [`BUANG_SENTRY_DSN`]                            | optional                                                                 | Sentry DSN                                                                                                                     |
| [`BUANG_TEMPORAL_ADDRESS`]                      | optional                                                                 | Temporal address. You can either use Temporal Cloud or self host it                                                            |
| [`BUANG_TEMPORAL_API_KEY`]                      | optional                                                                 | Temporal API key                                                                                                               |
| [`BUANG_TEMPORAL_NAMESPACE`]                    | optional                                                                 | Temporal namespace                                                                                                             |
| [`BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR`]            | defaults to `/buang/traefik`                                             | Traefik dynamic config directory                                                                                               |
| [`BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT`]     | defaults to `443`                                                        | Traefik HTTPS port                                                                                                             |
| [`BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT`]           | defaults to `80`                                                         | Traefik HTTP port                                                                                                              |
| [`GO_ENV`]                                      | defaults to `development`                                                | Golang environment                                                                                                             |

> [!TIP]
> If an environment variable is set to an empty value, `buang` behaves as if
> that variable is left undefined.

## `BUANG_API_KEY`

> Buang API key

The `BUANG_API_KEY` variable **MAY** be left undefined, in which case a default
value is used.

⚠️ This variable is **sensitive**; its value may contain private information.

## `BUANG_DBOS_ADMIN_SERVER_PORT`

> DBOS admin server port. Optional to use DBOS. Specify a port to enable the admin server, DBOS default is 3001

The `BUANG_DBOS_ADMIN_SERVER_PORT` variable **MAY** be left undefined.
Otherwise, the value **MUST** be a whole number.

```bash
export BUANG_DBOS_ADMIN_SERVER_PORT=-922337203685477888  # (non-normative)
export BUANG_DBOS_ADMIN_SERVER_PORT=+1844674407370954752 # (non-normative)
```

<details>
<summary>Signed integer syntax</summary>

Signed integers can only be specified using decimal notation. A leading positive
sign (`+`) is **OPTIONAL**. A leading negative sign (`-`) is **REQUIRED** in
order to specify a negative value.

Internally, the `BUANG_DBOS_ADMIN_SERVER_PORT` variable is represented using a
signed 64-bit integer type (`int`); any value that overflows this data-type is
invalid.

</details>

## `BUANG_DBOS_CONDUCTOR_API_KEY`

> DBOS conductor API key. Optional to use DBOS

The `BUANG_DBOS_CONDUCTOR_API_KEY` variable **MAY** be left undefined.

⚠️ This variable is **sensitive**; its value may contain private information.

## `BUANG_DBOS_CONDUCTOR_URL`

> DBOS conductor url. Optional to use DBOS. You can use either the DBOS console or self host it

The `BUANG_DBOS_CONDUCTOR_URL` variable **MAY** be left undefined.

```bash
export BUANG_DBOS_CONDUCTOR_URL=foo # (non-normative)
```

## `BUANG_DB_CONNECTION_STRING`

> Database connection string

The `BUANG_DB_CONNECTION_STRING` variable **MAY** be left undefined, in which
case a default value is used.

⚠️ This variable is **sensitive**; its value may contain private information.

## `BUANG_DB_TYPE`

> Database type

The `BUANG_DB_TYPE` variable **MAY** be left undefined, in which case the
default value of `sqlite` is used. Otherwise, the value **MUST** be either
`postgres` or `sqlite`.

```bash
export BUANG_DB_TYPE=postgres
export BUANG_DB_TYPE=sqlite   # (default)
```

## `BUANG_DURABLE_EXECUTOR`

> Durable executor

The `BUANG_DURABLE_EXECUTOR` variable **MAY** be left undefined, in which case
the default value of `temporal` is used. Otherwise, the value **MUST** be either
`temporal` or `dbos`.

```bash
export BUANG_DURABLE_EXECUTOR=temporal # (default)
export BUANG_DURABLE_EXECUTOR=dbos
```

## `BUANG_ENABLE_EXPERIMENTAL_BOOTSTRAP`

> Enable experimental bootstrap. Bootstrapping allows Buang to deploy itself and autoupdates on Saturdays at midnight every week

The `BUANG_ENABLE_EXPERIMENTAL_BOOTSTRAP` variable **MAY** be left undefined, in
which case the default value of `false` is used. Otherwise, the value **MUST**
be either `true` or `false`.

```bash
export BUANG_ENABLE_EXPERIMENTAL_BOOTSTRAP=true
export BUANG_ENABLE_EXPERIMENTAL_BOOTSTRAP=false # (default)
```

## `BUANG_ENABLE_TELEMETRY`

> Enable telemetry

The `BUANG_ENABLE_TELEMETRY` variable **MAY** be left undefined, in which case
the default value of `false` is used. Otherwise, the value **MUST** be either
`true` or `false`.

```bash
export BUANG_ENABLE_TELEMETRY=true
export BUANG_ENABLE_TELEMETRY=false # (default)
```

## `BUANG_EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE`

> Cron expression for when auto updates should run. Defaults to every Saturday at midnight

The `BUANG_EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE` variable **MAY** be left
undefined, in which case the default value of `0 0 * * 6` is used. Otherwise,
the value must be a valid 5 field cron expression.

```bash
export BUANG_EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE='0 0 * * 6' # (default)
```

## `BUANG_HOUSEKEEPING_PRUNE_DEPLOYMENTS`

> Cron expression for stale deployments should be pruned. Defaults to every 2 days

The `BUANG_HOUSEKEEPING_PRUNE_DEPLOYMENTS` variable **MAY** be left undefined,
in which case the default value of `0 0 */2 * *` is used. Otherwise, the value
must be a valid 5 field cron expression.

```bash
export BUANG_HOUSEKEEPING_PRUNE_DEPLOYMENTS='0 0 */2 * *' # (default)
```

## `BUANG_LOG_LEVEL`

> Log level

The `BUANG_LOG_LEVEL` variable **MAY** be left undefined, in which case the
default value of `INFO` is used. Otherwise, the value **MUST** be one of the
values shown in the examples below.

```bash
export BUANG_LOG_LEVEL=DEBUG
export BUANG_LOG_LEVEL=ERROR
export BUANG_LOG_LEVEL=INFO  # (default)
export BUANG_LOG_LEVEL=WARN
```

## `BUANG_OTEL_EXPORTER_OTLP_ENDPOINT`

> OpenTelemetry exporter endpoint

The `BUANG_OTEL_EXPORTER_OTLP_ENDPOINT` variable **MAY** be left undefined, in
which case the default value of `` is used.

```bash
export BUANG_OTEL_EXPORTER_OTLP_ENDPOINT='' # (default)
```

## `BUANG_OTEL_SERVICE_NAME`

> OpenTelemetry service name. This is also the DBOS application name which is required when using DBOS

The `BUANG_OTEL_SERVICE_NAME` variable **MAY** be left undefined, in which case
the default value of `skulpture-buang` is used.

```bash
export BUANG_OTEL_SERVICE_NAME=skulpture-buang # (default)
```

## `BUANG_SENTRY_DSN`

> Sentry DSN

The `BUANG_SENTRY_DSN` variable **MAY** be left undefined.

```bash
export BUANG_SENTRY_DSN=foo # (non-normative)
```

## `BUANG_TEMPORAL_ADDRESS`

> Temporal address. You can either use Temporal Cloud or self host it

The `BUANG_TEMPORAL_ADDRESS` variable **MAY** be left undefined.

```bash
export BUANG_TEMPORAL_ADDRESS=foo # (non-normative)
```

## `BUANG_TEMPORAL_API_KEY`

> Temporal API key

The `BUANG_TEMPORAL_API_KEY` variable **MAY** be left undefined.

⚠️ This variable is **sensitive**; its value may contain private information.

## `BUANG_TEMPORAL_NAMESPACE`

> Temporal namespace

The `BUANG_TEMPORAL_NAMESPACE` variable **MAY** be left undefined.

```bash
export BUANG_TEMPORAL_NAMESPACE=foo # (non-normative)
```

## `BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR`

> Traefik dynamic config directory

The `BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR` variable **MAY** be left undefined, in
which case the default value of `/buang/traefik` is used.

```bash
export BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR=/buang/traefik # (default)
```

## `BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT`

> Traefik HTTPS port

The `BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT` variable **MAY** be left
undefined, in which case the default value of `443` is used. Otherwise, the
value **MUST** be a valid network port.

```bash
export BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT=443   # (default)
export BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT=8000  # (non-normative) a port commonly used for private web servers
export BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT=https # (non-normative) the IANA service name that maps to port 443
```

<details>
<summary>Network port syntax</summary>

Ports may be specified as a numeric value no greater than `65535`.
Alternatively, a service name can be used. Service names are resolved against
the system's service database, typically located in the `/etc/service` file on
UNIX-like systems. Standard service names are published by IANA.

</details>

## `BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT`

> Traefik HTTP port

The `BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT` variable **MAY** be left undefined, in
which case the default value of `80` is used. Otherwise, the value **MUST** be a
valid network port.

```bash
export BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT=80    # (default)
export BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT=8000  # (non-normative) a port commonly used for private web servers
export BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT=https # (non-normative) the IANA service name that maps to port 443
```

<details>
<summary>Network port syntax</summary>

Ports may be specified as a numeric value no greater than `65535`.
Alternatively, a service name can be used. Service names are resolved against
the system's service database, typically located in the `/etc/service` file on
UNIX-like systems. Standard service names are published by IANA.

</details>

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

---

> [!NOTE]
> This document only describes environment variables declared using [Ferrite].
> `buang` may consume other undocumented environment variables.

> [!IMPORTANT]
> Some of the example values given in this document are **non-normative**.
> Although these values are syntactically valid, they may not be meaningful to
> `buang`.

<!-- references -->

[`buang_api_key`]: #BUANG_API_KEY
[`buang_db_connection_string`]: #BUANG_DB_CONNECTION_STRING
[`buang_db_type`]: #BUANG_DB_TYPE
[`buang_dbos_admin_server_port`]: #BUANG_DBOS_ADMIN_SERVER_PORT
[`buang_dbos_conductor_api_key`]: #BUANG_DBOS_CONDUCTOR_API_KEY
[`buang_dbos_conductor_url`]: #BUANG_DBOS_CONDUCTOR_URL
[`buang_durable_executor`]: #BUANG_DURABLE_EXECUTOR
[`buang_enable_experimental_bootstrap`]: #BUANG_ENABLE_EXPERIMENTAL_BOOTSTRAP
[`buang_enable_telemetry`]: #BUANG_ENABLE_TELEMETRY
[`buang_experimental_housekeeping_auto_update`]: #BUANG_EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE
[`buang_housekeeping_prune_deployments`]: #BUANG_HOUSEKEEPING_PRUNE_DEPLOYMENTS
[`buang_log_level`]: #BUANG_LOG_LEVEL
[`buang_otel_exporter_otlp_endpoint`]: #BUANG_OTEL_EXPORTER_OTLP_ENDPOINT
[`buang_otel_service_name`]: #BUANG_OTEL_SERVICE_NAME
[`buang_sentry_dsn`]: #BUANG_SENTRY_DSN
[`buang_temporal_address`]: #BUANG_TEMPORAL_ADDRESS
[`buang_temporal_api_key`]: #BUANG_TEMPORAL_API_KEY
[`buang_temporal_namespace`]: #BUANG_TEMPORAL_NAMESPACE
[`buang_traefik_dynamic_config_dir`]: #BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR
[`buang_traefik_entrypoint_web_port`]: #BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT
[`buang_traefik_entrypoint_websecure_port`]: #BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT
[ferrite]: https://github.com/dogmatiq/ferrite
[`go_env`]: #GO_ENV
