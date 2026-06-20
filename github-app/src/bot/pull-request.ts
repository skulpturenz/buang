import { type Probot } from "probot";
import { findRepoByFullName } from "../db/repos.js";
import { createDeployment, tearDownBranch } from "../services/buang.js";
import { decrypt } from "../services/crypto.js";
import {
	createGitHubDeployment,
	createSetupCheckRun,
	updateDeploymentStatus,
} from "../services/github.js";

export const registerPullRequestHandlers = (app: Probot) => {
	app.on(
		[
			"pull_request.opened",
			"pull_request.reopened",
			"pull_request.ready_for_review",
		],
		async ctx => {
			const { pull_request: pr, repository, installation } = ctx.payload;
			if (!installation) return;

			const repo = await findRepoByFullName(
				repository.full_name,
				installation.id,
			);

			const [owner, repoName] = repository.full_name.split("/");

			if (!repo?.buangProjectId) {
				await createSetupCheckRun(
					ctx.octokit,
					owner,
					repoName,
					pr.head.sha,
					repository.full_name,
					installation.id,
				);
				return;
			}

			if (repo.waitForWorkflow) return;

			const user = await getUserCreds(repo);
			if (!user) {
				app.log.error(
					{ repo: repository.full_name },
					"No user creds found for repo",
				);
				return;
			}

			const env = repo.envVarsEncrypted
				? (JSON.parse(decrypt(repo.envVarsEncrypted)) as Record<
						string,
						string
					>)
				: {};

			const deploymentId = await createGitHubDeployment(
				ctx.octokit,
				owner,
				repoName,
				pr.head.ref,
			);
			await updateDeploymentStatus(
				ctx.octokit,
				owner,
				repoName,
				deploymentId,
				"in_progress",
			);

			try {
				await createDeployment(user, repo.buangProjectId, {
					branch: pr.head.ref,
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
				app.log.error(err, "createDeployment failed");
				await updateDeploymentStatus(
					ctx.octokit,
					owner,
					repoName,
					deploymentId,
					"failure",
				);
			}
		},
	);

	app.on("pull_request.closed", async ctx => {
		const { pull_request: pr, repository, installation } = ctx.payload;
		if (!installation) return;

		const repo = await findRepoByFullName(
			repository.full_name,
			installation.id,
		);
		if (!repo?.buangProjectId) return;

		const user = await getUserCreds(repo);
		if (!user) return;

		try {
			await tearDownBranch(user, repo.buangProjectId, pr.head.ref);
		} catch (err) {
			app.log.error(err, "tearDownBranch failed");
		}

		const [owner, repoName] = repository.full_name.split("/");
		const { data: deployments } =
			await ctx.octokit.rest.repos.listDeployments({
				owner,
				repo: repoName,
				ref: pr.head.ref,
				environment: "preview",
			});
		for (const d of deployments) {
			await updateDeploymentStatus(
				ctx.octokit,
				owner,
				repoName,
				d.id,
				"inactive",
			);
		}
	});
};

const getUserCreds = async (repo: Awaited<ReturnType<typeof findRepoByFullName>>) => {
	if (!repo?.ownerUserId) return null;
	const { findUserById } = await import("../db/users.js");
	const user = await findUserById(repo.ownerUserId);
	if (!user?.buangApiBaseUrl || !user.buangApiKeyEncrypted) return null;
	return {
		baseUrl: user.buangApiBaseUrl,
		apiKey: decrypt(user.buangApiKeyEncrypted),
	};
};
