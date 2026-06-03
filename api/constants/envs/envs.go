package constantsenvs

import (
	"fmt"
	"log/slog"
	constantsfeaturetoggles "skulpture/buang/constants/feature_toggles"
	enumsdbtypes "skulpture/buang/enums/db_types"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"
	enumsenv "skulpture/buang/enums/env"
	"time"

	"github.com/adhocore/gronx"
	"github.com/dogmatiq/ferrite"
)

var (
	versionUnknown = "unknown"
	BUANG_VERSION  string
	GO_ENV         = ferrite.
			Enum("GO_ENV", "Golang environment").
			WithMembers(enumsenv.Production.String(), enumsenv.Development.String(), enumsenv.Test.String()).
			WithDefault(enumsenv.Development.String()).
			Required()
	BUANG_PORT = ferrite.NetworkPort("BUANG_PORT", "Port").
			WithDefault("80").
			Required()
	API_KEY = ferrite.
		String("BUANG_API_KEY", "Buang API key").
		WithSensitiveContent().
		WithDefault("supersecureapikey").
		Required()
	LOG_LEVEL = ferrite.EnumAs[slog.Level]("BUANG_LOG_LEVEL", "Log level").
			WithMembers(slog.LevelDebug, slog.LevelError, slog.LevelInfo, slog.LevelWarn).
			WithDefault(slog.LevelInfo).
			Required()
	ENABLE_TELEMETRY = ferrite.
				Bool("BUANG_ENABLE_TELEMETRY", "Enable telemetry").
				WithDefault(false).
				Required()
	OTEL_SERVICE_NAME = ferrite.
				String("BUANG_OTEL_SERVICE_NAME", "OpenTelemetry service name. This is also the DBOS application name which is required when using DBOS").
				WithDefault("skulpture-buang").
				Required()
	OTEL_EXPORTER_OTLP_ENDPOINT = ferrite.
					String("BUANG_OTEL_EXPORTER_OTLP_ENDPOINT", "OpenTelemetry exporter endpoint").
					WithDefault("").
					Optional()
	DB_TYPE = ferrite.
		Enum("BUANG_DB_TYPE", "Database type").
		WithMembers(enumsdbtypes.Pg.String(), enumsdbtypes.Sqlite.String()).
		WithDefault(enumsdbtypes.Sqlite.String()).
		Required()
	DB_CONNECTION_STRING = ferrite.
				String("BUANG_DB_CONNECTION_STRING", "Database connection string").
				WithSensitiveContent().
				WithDefault("file:test.db?_foreign_keys=true&mode=memory&cache=shared").
				Required()
	TEMPORAL_API_KEY = ferrite.
				String("BUANG_TEMPORAL_API_KEY", "Temporal API key").
				WithSensitiveContent().
				Optional()
	TEMPORAL_NAMESPACE = ferrite.
				String("BUANG_TEMPORAL_NAMESPACE", "Temporal namespace").
				Optional()
	TEMPORAL_ADDRESS = ferrite.
				String("BUANG_TEMPORAL_ADDRESS", "Temporal address. You can either use Temporal Cloud or self host it").
				Optional()
	DURABLE_EXECUTOR = ferrite.
				Enum("BUANG_DURABLE_EXECUTOR", "Durable executor").
				WithMembers(enumsdurableexecutors.Temporal.String(), enumsdurableexecutors.Dbos.String()).
				WithDefault(enumsdurableexecutors.Temporal.String()).
				Required()
	DBOS_CONDUCTOR_API_KEY = ferrite.
				String("BUANG_DBOS_CONDUCTOR_API_KEY", "DBOS conductor API key. Optional to use DBOS").
				WithSensitiveContent().
				Optional()
	DBOS_CONDUCTOR_URL = ferrite.
				String("BUANG_DBOS_CONDUCTOR_URL", "DBOS conductor url. Optional to use DBOS. You can use either the DBOS console or self host it").
				Optional()
	DBOS_ADMIN_SERVER_PORT = ferrite.
				Signed[int]("BUANG_DBOS_ADMIN_SERVER_PORT", "DBOS admin server port. Optional to use DBOS. Specify a port to enable the admin server, DBOS default is 3001").
				Optional()
	EXPERIMENTAL_BOOTSTRAP = ferrite.
				Bool(constantsfeaturetoggles.ENABLE_EXPERIMENTAL_BOOTSTRAP, "Enable experimental bootstrap. Bootstrapping allows Buang to deploy itself and autoupdates on Saturdays at midnight every week").
				WithDefault(false).
				Optional()
	HOUSEKEEPING_PRUNE_DEPLOYMENTS = ferrite.
					String("BUANG_HOUSEKEEPING_PRUNE_DEPLOYMENTS", "Cron expression for stale deployments should be pruned. Defaults to every 2 days").
					WithDefault("0 0 */2 * *").
					WithConstraint("must be a valid 5 field cron expression", gronx.IsValid).
					Required()
	EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE = ferrite.
						String("BUANG_EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE", "Cron expression for when auto updates should run. Defaults to every Saturday at midnight").
						WithDefault("0 0 * * 6").
						WithConstraint("must be a valid 5 field cron expression", gronx.IsValid).
						Optional()
	BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT = ferrite.
						NetworkPort("BUANG_TRAEFIK_ENTRYPOINT_WEB_PORT", "Traefik HTTP port").
						WithDefault("80").
						Optional()
	BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT = ferrite.
						NetworkPort("BUANG_TRAEFIK_ENTRYPOINT_WEBSECURE_PORT", "Traefik HTTPS port").
						WithDefault("443").
						Optional()
	BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR = ferrite.
						String("BUANG_TRAEFIK_DYNAMIC_CONFIG_DIR", "Traefik dynamic config directory").
						WithDefault("/buang/traefik").
						Required()
	BUANG_SENTRY_DSN = ferrite.
				String("BUANG_SENTRY_DSN", "Sentry DSN").
				Optional()
	INTERNAL_BUANG_BOOTSTRAP_PROJECT = "BUANG_BOOTSTRAP_PROJECT"
	INTERNAL_BUANG_BOOTSTRAP_DIR     = "BUANG_BOOTSTRAP_DIR"
	startTime                        = time.Now()
)

func init() {
	if v, _ := enumsenv.Parse(GO_ENV.Value()); v == enumsenv.Production && BUANG_VERSION == "" {
		panic(fmt.Sprintf("bad buang build, version: %v", versionUnknown))
	}

	ferrite.Init()
}

func Uptime() time.Duration {
	return time.Since(startTime)
}
