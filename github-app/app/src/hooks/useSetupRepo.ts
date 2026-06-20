import { useMutation, useQueryClient } from "@tanstack/react-query";
import { post } from "../lib/api.js";

export interface SetupRepoInput {
	repoFullName: string;
	installationId: number;
	composePath: string;
	serviceEntrypoint: string;
	requiresAuthn: boolean;
	username?: string;
	password?: string;
	envVars: Array<{ key: string; value: string }>;
	waitForWorkflow: boolean;
	workflowName?: string;
}

export const useSetupRepo = () => {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: (input: SetupRepoInput) => post("/repos/setup", input),
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["repo"] });
		},
	});
};
