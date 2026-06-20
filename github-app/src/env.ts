import { integer, string, url } from "austenite";
import { initialize } from "austenite/node";

export const DATABASE_URL = string(
	"BUANG_GHA_DATABASE_URL",
	"PostgreSQL connection string",
);

export const BUANG_GHA_SECRET_KEY = string(
	"BUANG_GHA_SECRET_KEY",
	"64-char hex string (32-byte AES-256 key)",
	{ isSensitive: true },
);

export const GITHUB_APP_ID = string(
	"BUANG_GHA_GITHUB_APP_ID",
	"GitHub App numeric ID",
);

export const GITHUB_PRIVATE_KEY = string(
	"BUANG_GHA_GITHUB_PRIVATE_KEY",
	"GitHub App private key in PEM format",
	{ isSensitive: true },
);

export const GITHUB_WEBHOOK_SECRET = string(
	"BUANG_GHA_GITHUB_WEBHOOK_SECRET",
	"GitHub webhook secret",
	{ isSensitive: true },
);

export const OIDC_ISSUER = url(
	"BUANG_GHA_OIDC_ISSUER",
	"OIDC provider issuer URL",
);

export const OIDC_CLIENT_ID = string(
	"BUANG_GHA_OIDC_CLIENT_ID",
	"OIDC client ID",
);

export const OIDC_CLIENT_SECRET = string(
	"BUANG_GHA_OIDC_CLIENT_SECRET",
	"OIDC client secret",
	{ isSensitive: true },
);

export const SESSION_SECRET = string(
	"BUANG_GHA_SESSION_SECRET",
	"Cookie session signing secret",
	{ isSensitive: true },
);

export const APP_BASE_URL = url(
	"BUANG_GHA_APP_BASE_URL",
	"Public base URL used in GitHub check run links",
);

export const PORT = integer("BUANG_GHA_PORT", "HTTP server port", {
	default: 3000,
});

export const SENTRY_DSN = string("BUANG_GHA_SENTRY_DSN", "Sentry DSN", {
	default: undefined,
});

initialize();
