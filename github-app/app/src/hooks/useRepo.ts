import { useQuery } from "@tanstack/react-query";
import { get } from "../lib/api.js";

export interface Repo {
	id: string;
	githubRepoFullName: string;
	isSetup: boolean;
	composePath: string | null;
	serviceEntrypoint: string | null;
	waitForWorkflow: boolean;
	workflowName: string | null;
	requiresAuthn: boolean;
}

export const useRepo = (repoFullName: string, installationId: string) =>
	useQuery<Repo>({
		queryKey: ["repo", repoFullName, installationId],
		queryFn: () =>
			get<Repo>(
				`/repos?repo=${encodeURIComponent(repoFullName)}&installation_id=${installationId}`,
			),
		enabled: !!repoFullName && !!installationId,
	});
