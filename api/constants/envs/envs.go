package constantsenvs

import (
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
	BUANG_VERSION = "unknown"
	GO_ENV        = ferrite.
			Enum("GO_ENV", "Golang environment").
			WithMembers(enumsenv.Production.String(), enumsenv.Development.String(), enumsenv.Test.String()).
			WithDefault(enumsenv.Development.String()).
			Required()
	API_KEY = ferrite.
		String("API_KEY", "Buang API key").
		WithSensitiveContent().
		WithDefault("supersecureapikey").
		Required()
	LOG_LEVEL = ferrite.EnumAs[slog.Level]("LOG_LEVEL", "Log level").
			WithMembers(slog.LevelDebug, slog.LevelError, slog.LevelInfo, slog.LevelWarn).
			WithDefault(slog.LevelInfo).
			Required()
	ENABLE_TELEMETRY = ferrite.
				Bool("ENABLE_TELEMETRY", "Enable telemetry").
				WithDefault(false).
				Required()
	OTEL_SERVICE_NAME = ferrite.
				String("OTEL_SERVICE_NAME", "OpenTelemetry service name. This is also the DBOS application name which is required when using DBOS").
				WithDefault("skulpture-buang").
				Required()
	OTEL_EXPORTER_OTLP_ENDPOINT = ferrite.
					String("OTEL_EXPORTER_OTLP_ENDPOINT", "OpenTelemetry exporter endpoint").
					WithDefault("").
					Optional()
	DB_TYPE = ferrite.
		Enum("DB_TYPE", "Database type").
		WithMembers(enumsdbtypes.Pg.String(), enumsdbtypes.Sqlite.String()).
		WithDefault(enumsdbtypes.Sqlite.String()).
		Required()
	DB_CONNECTION_STRING = ferrite.
				String("DB_CONNECTION_STRING", "Database connection string").
				WithSensitiveContent().
				WithDefault("file:test.db?_foreign_keys=true&mode=memory").
				Required()
	TEMPORAL_API_KEY = ferrite.
				String("TEMPORAL_API_KEY", "Temporal API key").
				WithSensitiveContent().
				Optional()
	TEMPORAL_NAMESPACE = ferrite.
				String("TEMPORAL_NAMESPACE", "Temporal namespace").
				Optional()
	TEMPORAL_ADDRESS = ferrite.
				String("TEMPORAL_ADDRESS", "Temporal address. You can either use Temporal Cloud or self host it").
				Optional()
	DURABLE_EXECUTOR = ferrite.
				Enum("DURABLE_EXECUTOR", "Durable executor").
				WithMembers(enumsdurableexecutors.Temporal.String(), enumsdurableexecutors.Dbos.String()).
				WithDefault(enumsdurableexecutors.Temporal.String()).
				Required()
	DBOS_CONDUCTOR_API_KEY = ferrite.
				String("DBOS_CONDUCTOR_API_KEY", "DBOS conductor API key. Optional to use DBOS").
				WithSensitiveContent().
				Optional()
	DBOS_CONDUCTOR_URL = ferrite.
				String("DBOS_CONDUCTOR_URL", "DBOS conductor url. Optional to use DBOS. You can use either the DBOS console or self host it").
				Optional()
	DBOS_ADMIN_SERVER_PORT = ferrite.
				Signed[int]("DBOS_ADMIN_SERVER_PORT", "DBOS admin server port. Optional to use DBOS. Specify a port to enable the admin server, DBOS default is 3001").
				Optional()
	EXPERIMENTAL_BOOTSTRAP = ferrite.
				Bool(constantsfeaturetoggles.EXPERIMENTAL_BOOTSTRAP, "Enable experimental bootstrap. Bootstrapping allows Buang to deploy itself and autoupdates on Saturdays at midnight every week").
				WithDefault(false).
				Optional()
	HOUSEKEEPING_PRUNE_DEPLOYMENTS = ferrite.
					String("HOUSEKEEPING_PRUNE_DEPLOYMENTS", "Cron expression for stale deployments should be pruned. Defaults to every 2 days").
					WithDefault("0 0 */2 * *").
					WithConstraint("must be a valid 5 field cron expression", gronx.IsValid).
					Required()
	EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE = ferrite.
						String("EXPERIMENTAL_HOUSEKEEPING_AUTO_UPDATE", "Cron expression for when auto updates should run. Defaults to every Saturday at midnight").
						WithDefault("0 0 * * 6").
						WithConstraint("must be a valid 5 field cron expression", gronx.IsValid).
						Optional()
	startTime = time.Now()
)

func init() {
	ferrite.Init()
}

func Uptime() time.Duration {
	return time.Since(startTime)
}
