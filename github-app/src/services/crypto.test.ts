import { beforeAll, describe, expect, it } from "vitest";

// Set required env vars before the module initializes
beforeAll(() => {
	process.env.BUANG_GHA_DATABASE_URL = "postgres://test:test@localhost:5432/test";
	process.env.BUANG_GHA_SECRET_KEY =
		"0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20";
	process.env.BUANG_GHA_GITHUB_APP_ID = "1";
	process.env.BUANG_GHA_GITHUB_PRIVATE_KEY = "dummy";
	process.env.BUANG_GHA_GITHUB_WEBHOOK_SECRET = "dummy";
	process.env.BUANG_GHA_OIDC_ISSUER = "https://example.com";
	process.env.BUANG_GHA_OIDC_CLIENT_ID = "dummy";
	process.env.BUANG_GHA_OIDC_CLIENT_SECRET = "dummy";
	process.env.BUANG_GHA_SESSION_SECRET = "dummy";
	process.env.BUANG_GHA_APP_BASE_URL = "http://localhost:3000";
});

describe("crypto", () => {
	it("round-trips plaintext strings", async () => {
		const { encrypt, decrypt } = await import("./crypto.js");
		const original = "super-secret-api-key";
		const ciphertext = encrypt(original);
		expect(ciphertext).not.toBe(original);
		expect(decrypt(ciphertext)).toBe(original);
	});

	it("produces unique ciphertexts for the same plaintext", async () => {
		const { encrypt } = await import("./crypto.js");
		const a = encrypt("same-value");
		const b = encrypt("same-value");
		expect(a).not.toBe(b);
	});

	it("round-trips JSON env var maps", async () => {
		const { encrypt, decrypt } = await import("./crypto.js");
		const envVars = { FOO: "bar", DB_URL: "postgres://localhost/db" };
		const ciphertext = encrypt(JSON.stringify(envVars));
		const decoded = JSON.parse(decrypt(ciphertext));
		expect(decoded).toEqual(envVars);
	});

	it("throws on tampered ciphertext", async () => {
		const { encrypt, decrypt } = await import("./crypto.js");
		const ciphertext = encrypt("value");
		const buf = Buffer.from(ciphertext, "base64");
		// Flip a byte in the ciphertext body
		buf[20] ^= 0xff;
		expect(() => decrypt(buf.toString("base64"))).toThrow();
	});
});
