import { Hono } from "hono";
import {
	type InferType,
	ValidationError,
	object,
	string as yupString,
} from "yup";
import { findUserById, updateUser } from "../db/users.js";
import { encrypt } from "../services/crypto.js";
import "./context.js";
import { requireAuth } from "./auth.js";

export const usersRouter = new Hono();

usersRouter.use("*", requireAuth);

usersRouter.get("/me", async c => {
	const userId: bigint = c.get("userId");
	const user = await findUserById(userId);
	if (!user) return c.json({ errors: ["Not found"] }, 404);
	return c.json({
		id: String(user.id),
		email: user.email,
		authType: user.authType,
		buangApiBaseUrl: user.buangApiBaseUrl,
		hasBuangKey: !!user.buangApiKeyEncrypted,
	});
});

const patchSchema = object({
	buangApiBaseUrl: yupString().url(),
	buangApiKey: yupString(),
});

usersRouter.patch("/me", async c => {
	let body: InferType<typeof patchSchema>;
	try {
		body = await patchSchema.validate(await c.req.json(), {
			abortEarly: false,
		});
	} catch (err) {
		if (err instanceof ValidationError) {
			return c.json({ errors: err.errors }, 400);
		}
		throw err;
	}

	const userId: bigint = c.get("userId");
	const updated = await updateUser(userId, {
		buangApiBaseUrl: body.buangApiBaseUrl,
		buangApiKeyEncrypted: body.buangApiKey
			? encrypt(body.buangApiKey)
			: undefined,
	});

	return c.json({
		id: String(updated.id),
		buangApiBaseUrl: updated.buangApiBaseUrl,
		hasBuangKey: !!updated.buangApiKeyEncrypted,
	});
});
