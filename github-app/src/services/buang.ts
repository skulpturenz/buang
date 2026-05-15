import { type BuangCreds, type CreateDeploymentParams } from "../types.js";

export interface CreateProjectParams {
	repository: string;
	composePath: string;
	requiresAuthn: boolean;
	username?: string;
	password?: string;
}

export const createProject = async (
	creds: BuangCreds,
	params: CreateProjectParams,
): Promise<string> => {
	const res = await fetch(`${creds.baseUrl}/api/v1/project`, {
		method: "POST",
		headers: {
			"X-API-Key": creds.apiKey,
			"Content-Type": "application/json",
		},
		body: JSON.stringify(params),
	});
	if (!res.ok) {
		const body = await res.text();
		throw new Error(`createProject failed (${res.status}): ${body}`);
	}
	return res.text();
};

export const createDeployment = async (
	creds: BuangCreds,
	projectId: bigint,
	params: CreateDeploymentParams,
): Promise<string> => {
	const res = await fetch(
		`${creds.baseUrl}/api/v1/project/${projectId}/deployment?waitForDeployment=true`,
		{
			method: "POST",
			headers: {
				"X-API-Key": creds.apiKey,
				"Content-Type": "application/json",
			},
			body: JSON.stringify(params),
		},
	);
	if (!res.ok) {
		const body = await res.text();
		throw new Error(`createDeployment failed (${res.status}): ${body}`);
	}
	return res.text();
};

export const tearDownBranch = async (
	creds: BuangCreds,
	projectId: bigint,
	branch: string,
): Promise<void> => {
	const res = await fetch(
		`${creds.baseUrl}/api/v1/project/${projectId}/branch`,
		{
			method: "DELETE",
			headers: {
				"X-API-Key": creds.apiKey,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ branch }),
		},
	);
	if (!res.ok && res.status !== 404) {
		throw new Error(`tearDownBranch failed (${res.status})`);
	}
};
