import { Hono } from "hono";
import {
	type InferType,
	ValidationError,
	array,
	boolean as yupBoolean,
	number as yupNumber,
	object,
	string as yupString,
} from "yup";
import { findRepoByFullName, updateRepo, upsertRepo } from "../db/repos.js";
import { findUserById } from "../db/users.js";
import { createProject } from "../services/buang.js";
import { encrypt } from "../services/crypto.js";
import "./context.js";
import { requireAuth } from "./auth.js";

export const reposRouter = new Hono();

reposRouter.use("*", requireAuth);

reposRouter.get("/repos", async c => {
	const repoFullName = c.req.query("repo");
	const rawInstallationId = c.req.query("installation_id");
	if (!repoFullName || !rawInstallationId) {
		return c.json({ errors: ["repo and installation_id are required"] }, 400);
	}
	const repo = await findRepoByFullName(
		repoFullName,
		Number(rawInstallationId),
	);
	if (!repo) return c.json({ errors: ["Not found"] }, 404);
	return c.json({
		id: String(repo.id),
		githubRepoFullName: repo.githubRepoFullName,
		isSetup: !!repo.buangProjectId,
		composePath: repo.composePath,
		serviceEntrypoint: repo.serviceEntrypoint,
		waitForWorkflow: repo.waitForWorkflow,
		workflowName: repo.workflowName,
		requiresAuthn: repo.requiresAuthn,
	});
});

const setupSchema = object({
	repoFullName: yupString().required(),
	installationId: yupNumber().required(),
	composePath: yupString().required(),
	serviceEntrypoint: yupString().required(),
	requiresAuthn: yupBoolean().required(),
	username: yupString().when("requiresAuthn", {
		is: true,
		then: s => s.required(),
	}),
	password: yupString().when("requiresAuthn", {
		is: true,
		then: s => s.required(),
	}),
	envVars: array(
		object({
			key: yupString().required(),
			value: yupString().required(),
		}),
	).default([]),
	waitForWorkflow: yupBoolean().default(false),
	workflowName: yupString().nullable(),
});

reposRouter.post("/repos/setup", async c => {
	let body: InferType<typeof setupSchema>;
	try {
		body = await setupSchema.validate(await c.req.json(), {
			abortEarly: false,
		});
	} catch (err) {
		if (err instanceof ValidationError) {
			return c.json({ errors: err.errors }, 400);
		}
		throw err;
	}

	const userId: bigint = c.get("userId");
	const user = await findUserById(userId);
	if (!user?.buangApiBaseUrl || !user.buangApiKeyEncrypted) {
		return c.json(
			{ errors: ["Complete your profile before setting up a repository"] },
			422,
		);
	}

	const { decrypt } = await import("../services/crypto.js");
	const creds = {
		baseUrl: user.buangApiBaseUrl,
		apiKey: decrypt(user.buangApiKeyEncrypted),
	};

	const projectIdStr = await createProject(creds, {
		repository: body.repoFullName!,
		composePath: body.composePath!,
		requiresAuthn: body.requiresAuthn!,
		username: body.username,
		password: body.password,
	});

	const envMap: Record<string, string> = {};
	for (const { key, value } of body.envVars ?? []) {
		if (key && value) envMap[key] = value;
	}

	const repo = await upsertRepo(body.repoFullName!, body.installationId!);
	await updateRepo(repo.id, {
		buangProjectId: BigInt(projectIdStr.trim()),
		ownerUserId: userId,
		envVarsEncrypted:
			Object.keys(envMap).length > 0
				? encrypt(JSON.stringify(envMap))
				: undefined,
		composePath: body.composePath,
		serviceEntrypoint: body.serviceEntrypoint,
		waitForWorkflow: body.waitForWorkflow,
		workflowName: body.workflowName ?? null,
		requiresAuthn: body.requiresAuthn,
	});

	return c.json({ ok: true, projectId: projectIdStr.trim() }, 201);
});
