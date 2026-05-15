import bcrypt from "bcrypt";
import {
	getAuth,
	oidcAuthMiddleware,
	processOAuthCallback,
} from "@hono/oidc-auth";
import { Hono } from "hono";
import { getCookie, setCookie } from "hono/cookie";
import {
	type InferType,
	ValidationError,
	object,
	string as yupString,
} from "yup";
import { createUser, findUserByEmail } from "../db/users.js";
import { encrypt } from "../services/crypto.js";
import "./context.js";

export const authRouter = new Hono();

const registerSchema = object({
	email: yupString().email().required(),
	password: yupString().min(8).required(),
	buangApiBaseUrl: yupString().url().required(),
	buangApiKey: yupString().required(),
});

const loginSchema = object({
	email: yupString().email().required(),
	password: yupString().required(),
});

authRouter.post("/register", async c => {
	let body: InferType<typeof registerSchema>;
	try {
		body = await registerSchema.validate(await c.req.json(), {
			abortEarly: false,
		});
	} catch (err) {
		if (err instanceof ValidationError) {
			return c.json({ errors: err.errors }, 400);
		}
		throw err;
	}

	const existing = await findUserByEmail(body.email);
	if (existing) {
		return c.json({ errors: ["Email already registered"] }, 409);
	}

	// Verify the Buang API is reachable
	try {
		const ping = await fetch(body.buangApiBaseUrl, {
			method: "HEAD",
			signal: AbortSignal.timeout(5000),
		});
		if (!ping.ok && ping.status >= 500) {
			return c.json(
				{ errors: ["Buang API URL is not reachable"] },
				422,
			);
		}
	} catch (_err) {
		return c.json({ errors: ["Buang API URL is not reachable"] }, 422);
	}

	const [passwordHash, buangApiKeyEncrypted] = await Promise.all([
		bcrypt.hash(body.password, 12),
		Promise.resolve(encrypt(body.buangApiKey)),
	]);

	const user = await createUser({
		email: body.email,
		passwordHash,
		authType: "email",
		buangApiBaseUrl: body.buangApiBaseUrl,
		buangApiKeyEncrypted,
	});

	setCookie(c, "session_user_id", String(user.id), {
		httpOnly: true,
		sameSite: "Lax",
		path: "/",
	});

	return c.json({ id: String(user.id) }, 201);
});

authRouter.post("/login", async c => {
	let body: InferType<typeof loginSchema>;
	try {
		body = await loginSchema.validate(await c.req.json(), {
			abortEarly: false,
		});
	} catch (err) {
		if (err instanceof ValidationError) {
			return c.json({ errors: err.errors }, 400);
		}
		throw err;
	}

	const user = await findUserByEmail(body.email);
	if (!user || !user.passwordHash) {
		return c.json({ errors: ["Invalid credentials"] }, 401);
	}

	const valid = await bcrypt.compare(body.password, user.passwordHash);
	if (!valid) {
		return c.json({ errors: ["Invalid credentials"] }, 401);
	}

	setCookie(c, "session_user_id", String(user.id), {
		httpOnly: true,
		sameSite: "Lax",
		path: "/",
	});

	return c.json({ id: String(user.id) });
});

authRouter.post("/logout", c => {
	setCookie(c, "session_user_id", "", { maxAge: 0, path: "/" });
	return c.json({ ok: true });
});

// OIDC login — browser redirect triggers the IdP authorization flow
authRouter.get("/auth/oidc", oidcAuthMiddleware(), async c => {
	// If the middleware passes (user is already authenticated), redirect to setup
	return c.redirect("/setup");
});

// OIDC callback — IdP redirects here after successful authentication
authRouter.get("/auth/oidc/callback", async c => {
	// Complete the OIDC token exchange; the library sets its own session cookie
	await processOAuthCallback(c);

	const auth = await getAuth(c);
	if (!auth) return c.text("OIDC authentication failed", 401);

	const email = auth.email ?? auth.sub ?? null;

	// Upsert user — find by email if we have one, else create new
	let user = email ? await findUserByEmail(email) : null;
	if (!user) {
		user = await createUser({ email, passwordHash: null, authType: "oauth" });
	}

	// Set our own session cookie so requireAuth works uniformly
	setCookie(c, "session_user_id", String(user.id), {
		httpOnly: true,
		sameSite: "Lax",
		path: "/",
	});

	// If user hasn't provided Buang credentials yet, send them to complete profile
	const dest =
		user.buangApiKeyEncrypted ? "/setup" : "/complete-profile";
	return c.redirect(dest);
});

export const requireAuth = async (c: any, next: () => Promise<void>) => {
	const userId = getCookie(c, "session_user_id");
	if (!userId) return c.json({ errors: ["Unauthorized"] }, 401);
	c.set("userId", BigInt(userId));
	await next();
};
