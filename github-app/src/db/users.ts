import { type AuthType, type User } from "../types.js";
import { pool } from "./client.js";

const rowToUser = (row: any): User => ({
	id: row.id,
	email: row.email,
	passwordHash: row.password_hash,
	authType: row.auth_type as AuthType,
	buangApiBaseUrl: row.buang_api_base_url,
	buangApiKeyEncrypted: row.buang_api_key_encrypted,
	createdAt: row.created_at,
	updatedAt: row.updated_at,
});

export const findUserByEmail = async (email: string): Promise<User | null> => {
	const { rows } = await pool.query(
		"SELECT * FROM users WHERE email = $1 LIMIT 1",
		[email],
	);
	return rows[0] ? rowToUser(rows[0]) : null;
};

export const findUserById = async (id: bigint): Promise<User | null> => {
	const { rows } = await pool.query(
		"SELECT * FROM users WHERE id = $1 LIMIT 1",
		[id],
	);
	return rows[0] ? rowToUser(rows[0]) : null;
};

export interface CreateUserInput {
	email: string | null;
	passwordHash: string | null;
	authType: AuthType;
	buangApiBaseUrl?: string;
	buangApiKeyEncrypted?: string;
}

export const createUser = async (input: CreateUserInput): Promise<User> => {
	const { rows } = await pool.query(
		`INSERT INTO users
			(email, password_hash, auth_type, buang_api_base_url, buang_api_key_encrypted)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`,
		[
			input.email,
			input.passwordHash,
			input.authType,
			input.buangApiBaseUrl ?? null,
			input.buangApiKeyEncrypted ?? null,
		],
	);
	return rowToUser(rows[0]);
};

export interface UpdateUserInput {
	buangApiBaseUrl?: string;
	buangApiKeyEncrypted?: string;
}

export const updateUser = async (
	id: bigint,
	input: UpdateUserInput,
): Promise<User> => {
	const { rows } = await pool.query(
		`UPDATE users
		SET buang_api_base_url = COALESCE($2, buang_api_base_url),
		    buang_api_key_encrypted = COALESCE($3, buang_api_key_encrypted)
		WHERE id = $1
		RETURNING *`,
		[id, input.buangApiBaseUrl ?? null, input.buangApiKeyEncrypted ?? null],
	);
	return rowToUser(rows[0]);
};
