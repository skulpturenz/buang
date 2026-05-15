import { readFileSync } from "fs";
import { createServer } from "http";
import { serve } from "@hono/node-server";
import { serveStatic } from "@hono/node-server/serve-static";
import { init as sentryInit } from "@sentry/node";
import { Hono } from "hono";
import { createNodeMiddleware, createProbot } from "probot";
import { apiRouter } from "./api/index.js";
import { probotApp } from "./bot/index.js";
import { PORT, SENTRY_DSN } from "./env.js";

const sentryDsn = SENTRY_DSN.value();
if (sentryDsn) {
	sentryInit({ dsn: sentryDsn });
}

const probot = createProbot();
const webhookMiddleware = createNodeMiddleware(probotApp, {
	probot,
	path: "/api/github/webhooks",
});

const app = new Hono<{ Variables: { userId: bigint } }>();

// Bridge Probot's Node.js middleware into Hono
app.post("/api/github/webhooks", async c => {
	await new Promise<void>((resolve, reject) => {
		const nodeReq = (c.env as any).incoming;
		const nodeRes = (c.env as any).outgoing;
		webhookMiddleware(nodeReq, nodeRes, (err?: any) => {
			if (err) reject(err as Error);
			else resolve();
		});
	});
	return new Response(null, { status: 200 });
});

app.route("/api", apiRouter);

app.use("/assets/*", serveStatic({ root: "./app/dist" }));

app.get("*", c => {
	try {
		const html = readFileSync("./app/dist/index.html", "utf8");
		return c.html(html);
	} catch (_err) {
		return c.text("Frontend not built. Run `pnpm build:app` first.", 503);
	}
});

const port = PORT.value();

serve(
	{ fetch: app.fetch, port, createServer },
	() => { console.log(`GitHub App listening on port ${port}`); },
);
