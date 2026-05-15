import { type Probot } from "probot";
import { findRepoByFullName } from "../db/repos.js";
import { createDeployment } from "../services/buang.js";
import { decrypt } from "../services/crypto.js";
import {
	createGitHubDeployment,
	findOpenPRForBranch,
	updateDeploymentStatus,
} from "../services/github.js";

export const registerWorkflowRunHandlers = (app: Probot) => {
	app.on("workflow_run.completed", async ctx => {
		const { workflow_run: run, repository, installation } = ctx.payload;
		if (!installation) return;
		if (run.conclusion !== "success") return;

		const repo = await findRepoByFullName(
			repository.full_name,
			installation.id,
		);
		if (!repo?.buangProjectId) return;
		if (!repo.waitForWorkflow) return;
		if (repo.workflowName && repo.workflowName !== run.name) return;

		const [owner, repoName] = repository.full_name.split("/");
		const pr = await findOpenPRForBranch(
			ctx.octokit,
			owner,
			repoName,
			run.head_branch,
		);
		if (!pr) return;

		if (!repo.ownerUserId) return;
		const { findUserById } = await import("../db/users.js");
		const user = await findUserById(repo.ownerUserId);
		if (!user?.buangApiBaseUrl || !user.buangApiKeyEncrypted) return;

		const creds = {
			baseUrl: user.buangApiBaseUrl,
			apiKey: decrypt(user.buangApiKeyEncrypted),
		};
		const env = repo.envVarsEncrypted
			? (JSON.parse(decrypt(repo.envVarsEncrypted)) as Record<string, string>)
			: {};

		const deploymentId = await createGitHubDeployment(
			ctx.octokit,
			owner,
			repoName,
			run.head_branch,
		);
		await updateDeploymentStatus(
			ctx.octokit,
			owner,
			repoName,
			deploymentId,
			"in_progress",
		);

		try {
			await createDeployment(creds, repo.buangProjectId, {
				branch: run.head_branch,
				sha: pr.head.sha,
				serviceEntrypoint: repo.serviceEntrypoint ?? "traefik:80",
				env,
			});
			await updateDeploymentStatus(
				ctx.octokit,
				owner,
				repoName,
				deploymentId,
				"success",
			);
		} catch (err) {
			app.log.error(err, "createDeployment via workflow_run failed");
			await updateDeploymentStatus(
				ctx.octokit,
				owner,
				repoName,
				deploymentId,
				"failure",
			);
		}
	});
};
