import { type Probot } from "probot";
import { upsertRepo } from "../db/repos.js";
import { createSetupCheckRun } from "../services/github.js";

export const registerInstallationHandlers = (app: Probot) => {
	const handleRepos = async (
		repos: Array<{ full_name: string }>,
		installationId: number,
		octokit: any,
	) => {
		for (const r of repos) {
			await upsertRepo(r.full_name, installationId);
			const [owner, repo] = r.full_name.split("/");
			// We need a sha to create a check run — use the default branch HEAD
			try {
				const { data: branch } =
					await octokit.rest.repos.getBranch({
						owner,
						repo,
						branch: "HEAD",
					});
				await createSetupCheckRun(
					octokit,
					owner,
					repo,
					branch.commit.sha,
					r.full_name,
					installationId,
				);
			} catch (_err) {
				// Skip if we can't reach the repo (e.g. empty repo)
				app.log.warn(
					{ repo: r.full_name },
					"Could not create setup check run",
				);
			}
		}
	};

	app.on("installation.created", async ctx => {
		const { installation, repositories } = ctx.payload;
		await handleRepos(
			repositories ?? [],
			installation.id,
			ctx.octokit,
		);
	});

	app.on("installation_repositories.added", async ctx => {
		const { installation, repositories_added } = ctx.payload;
		await handleRepos(
			repositories_added,
			installation.id,
			ctx.octokit,
		);
	});
};
