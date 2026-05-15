export type AuthType = "email" | "oauth";

export interface User {
	id: bigint;
	email: string | null;
	passwordHash: string | null;
	authType: AuthType;
	buangApiBaseUrl: string | null;
	buangApiKeyEncrypted: string | null;
	createdAt: Date;
	updatedAt: Date;
}

export interface Repository {
	id: bigint;
	githubRepoFullName: string;
	installationId: bigint;
	buangProjectId: bigint | null;
	ownerUserId: bigint | null;
	envVarsEncrypted: string | null;
	composePath: string | null;
	serviceEntrypoint: string | null;
	waitForWorkflow: boolean;
	workflowName: string | null;
	requiresAuthn: boolean;
	createdAt: Date;
	updatedAt: Date;
}

export interface BuangCreds {
	baseUrl: string;
	apiKey: string;
}

export interface CreateDeploymentParams {
	branch: string;
	sha: string;
	serviceEntrypoint: string;
	env: Record<string, string>;
}
