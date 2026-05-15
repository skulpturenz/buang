import { type Repository } from "../types.js";
import { pool } from "./client.js";

const rowToRepo = (row: any): Repository => ({
	id: row.id,
	githubRepoFullName: row.github_repo_full_name,
	installationId: row.installation_id,
	buangProjectId: row.buang_project_id,
	ownerUserId: row.owner_user_id,
	envVarsEncrypted: row.env_vars_encrypted,
	composePath: row.compose_path,
	serviceEntrypoint: row.service_entrypoint,
	waitForWorkflow: row.wait_for_workflow,
	workflowName: row.workflow_name,
	requiresAuthn: row.requires_authn,
	createdAt: row.created_at,
	updatedAt: row.updated_at,
});

export const findRepoByFullName = async (
	fullName: string,
	installationId: number,
): Promise<Repository | null> => {
	const { rows } = await pool.query(
		"SELECT * FROM repositories WHERE github_repo_full_name = $1 AND installation_id = $2 LIMIT 1",
		[fullName, installationId],
	);
	return rows[0] ? rowToRepo(rows[0]) : null;
};

export const findRepoById = async (id: bigint): Promise<Repository | null> => {
	const { rows } = await pool.query(
		"SELECT * FROM repositories WHERE id = $1 LIMIT 1",
		[id],
	);
	return rows[0] ? rowToRepo(rows[0]) : null;
};

export const upsertRepo = async (
	fullName: string,
	installationId: number,
): Promise<Repository> => {
	const { rows } = await pool.query(
		`INSERT INTO repositories (github_repo_full_name, installation_id)
		VALUES ($1, $2)
		ON CONFLICT (github_repo_full_name, installation_id) DO UPDATE
		SET updated_at = NOW()
		RETURNING *`,
		[fullName, installationId],
	);
	return rowToRepo(rows[0]);
};

export interface UpdateRepoInput {
	buangProjectId?: bigint;
	ownerUserId?: bigint;
	envVarsEncrypted?: string;
	composePath?: string;
	serviceEntrypoint?: string;
	waitForWorkflow?: boolean;
	workflowName?: string | null;
	requiresAuthn?: boolean;
}

export const updateRepo = async (
	id: bigint,
	input: UpdateRepoInput,
): Promise<Repository> => {
	const { rows } = await pool.query(
		`UPDATE repositories SET
			buang_project_id      = COALESCE($2, buang_project_id),
			owner_user_id         = COALESCE($3, owner_user_id),
			env_vars_encrypted    = COALESCE($4, env_vars_encrypted),
			compose_path          = COALESCE($5, compose_path),
			service_entrypoint    = COALESCE($6, service_entrypoint),
			wait_for_workflow     = COALESCE($7, wait_for_workflow),
			workflow_name         = COALESCE($8, workflow_name),
			requires_authn        = COALESCE($9, requires_authn)
		WHERE id = $1
		RETURNING *`,
		[
			id,
			input.buangProjectId ?? null,
			input.ownerUserId ?? null,
			input.envVarsEncrypted ?? null,
			input.composePath ?? null,
			input.serviceEntrypoint ?? null,
			input.waitForWorkflow ?? null,
			input.workflowName !== undefined ? input.workflowName : null,
			input.requiresAuthn ?? null,
		],
	);
	return rowToRepo(rows[0]);
};
