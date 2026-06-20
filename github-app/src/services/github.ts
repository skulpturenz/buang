import { type ProbotOctokit } from "probot";
import { APP_BASE_URL } from "../env.js";

type Octokit = InstanceType<typeof ProbotOctokit>;

export const createSetupCheckRun = async (
	octokit: Octokit,
	owner: string,
	repo: string,
	sha: string,
	repoFullName: string,
	installationId: number,
	update = false,
): Promise<void> => {
	const params = new URLSearchParams({
		repo: repoFullName,
		installation_id: String(installationId),
	});
	if (update) params.set("update", "true");
	const setupUrl = `${APP_BASE_URL.value().href.replace(/\/$/, "")}/setup?${params}`;

	await octokit.rest.checks.create({
		owner,
		repo,
		name: "Buang Setup",
		head_sha: sha,
		status: "completed",
		conclusion: "action_required",
		output: {
			title: update
				? "Repository authentication update required"
				: "Repository setup required",
			summary: update
				? "This repository was made private. Please update the authentication settings for the Buang deployment to continue."
				: "Please configure this repository with Buang to enable preview deployments.",
		},
		details_url: setupUrl,
	});
};

export const createGitHubDeployment = async (
	octokit: Octokit,
	owner: string,
	repo: string,
	ref: string,
): Promise<number> => {
	const { data } = await octokit.rest.repos.createDeployment({
		owner,
		repo,
		ref,
		environment: "preview",
		auto_merge: false,
		required_contexts: [],
		description: "Buang preview deployment",
	});
	if (!("id" in data)) throw new Error("Deployment creation returned no ID");
	return data.id;
};

export const updateDeploymentStatus = async (
	octokit: Octokit,
	owner: string,
	repo: string,
	deploymentId: number,
	state: "success" | "failure" | "inactive" | "in_progress",
	environmentUrl?: string,
): Promise<void> => {
	await octokit.rest.repos.createDeploymentStatus({
		owner,
		repo,
		deployment_id: deploymentId,
		state,
		environment_url: environmentUrl,
		description:
			state === "success"
				? "Preview deployment is live"
				: state === "inactive"
					? "Preview deployment torn down"
					: state === "in_progress"
						? "Deploying preview…"
						: "Deployment failed",
	});
};

export const findOpenPRForBranch = async (
	octokit: Octokit,
	owner: string,
	repo: string,
	branch: string,
): Promise<{ number: number; head: { sha: string } } | null> => {
	const { data } = await octokit.rest.pulls.list({
		owner,
		repo,
		state: "open",
		head: `${owner}:${branch}`,
	});
	return data[0] ?? null;
};
