import { type Probot } from "probot";
import { findRepoByFullName } from "../db/repos.js";
import { createSetupCheckRun } from "../services/github.js";

export const registerRepositoryHandlers = (app: Probot) => {
	app.on("repository.privatized", async ctx => {
		const { repository, installation } = ctx.payload;
		if (!installation) return;

		const repo = await findRepoByFullName(
			repository.full_name,
			installation.id,
		);
		if (!repo) return;

		const [owner, repoName] = repository.full_name.split("/");
		try {
			const { data: branch } = await ctx.octokit.rest.repos.getBranch({
				owner,
				repo: repoName,
				branch: "HEAD",
			});
			await createSetupCheckRun(
				ctx.octokit,
				owner,
				repoName,
				branch.commit.sha,
				repository.full_name,
				installation.id,
				true,
			);
		} catch (_err) {
			app.log.warn(
				{ repo: repository.full_name },
				"Could not create update check run",
			);
		}
	});
};
