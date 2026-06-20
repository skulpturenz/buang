import { initOidcAuthMiddleware } from "@hono/oidc-auth";
import { Hono } from "hono";
import { logger } from "hono/logger";
import { requestId } from "hono/request-id";
import { secureHeaders } from "hono/secure-headers";
import { timing } from "hono/timing";
import {
	APP_BASE_URL,
	OIDC_CLIENT_ID,
	OIDC_CLIENT_SECRET,
	OIDC_ISSUER,
	SESSION_SECRET,
} from "../env.js";
import { authRouter } from "./auth.js";
import { reposRouter } from "./repos.js";
import { usersRouter } from "./users.js";

export const apiRouter = new Hono();

apiRouter.use("*", secureHeaders());
apiRouter.use("*", requestId());
apiRouter.use("*", logger());
apiRouter.use("*", timing());

// Configure @hono/oidc-auth with our austenite env vars
apiRouter.use(
	"*",
	initOidcAuthMiddleware({
		OIDC_AUTH_SECRET: SESSION_SECRET.value(),
		OIDC_ISSUER: OIDC_ISSUER.value().href,
		OIDC_CLIENT_ID: OIDC_CLIENT_ID.value(),
		OIDC_CLIENT_SECRET: OIDC_CLIENT_SECRET.value(),
		OIDC_REDIRECT_URI: `${APP_BASE_URL.value().href.replace(/\/$/, "")}/api/auth/oidc/callback`,
		OIDC_SCOPES: "openid email profile",
	}),
);

apiRouter.route("/", authRouter);
apiRouter.route("/", usersRouter);
apiRouter.route("/", reposRouter);
