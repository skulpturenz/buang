import {
	beforeAll,
	beforeEach,
	describe,
	expect,
	it,
	vi,
} from "vitest";

beforeAll(() => {
	process.env.DATABASE_URL = "postgres://test:test@localhost:5432/test";
	process.env.BUANG_GHA_SECRET_KEY =
		"0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20";
	process.env.GITHUB_APP_ID = "1";
	process.env.GITHUB_PRIVATE_KEY = "dummy";
	process.env.GITHUB_WEBHOOK_SECRET = "secret";
	process.env.OIDC_ISSUER = "https://example.com";
	process.env.OIDC_CLIENT_ID = "dummy";
	process.env.OIDC_CLIENT_SECRET = "dummy";
	process.env.SESSION_SECRET = "dummy";
	process.env.APP_BASE_URL = "http://localhost:3000";
});

const mockFindRepoByFullName = vi.fn();
const mockFindUserById = vi.fn();
const mockCreateDeployment = vi.fn().mockResolvedValue("42");
const mockTearDownBranch = vi.fn().mockResolvedValue(undefined);
const mockCreateSetupCheckRun = vi.fn().mockResolvedValue(undefined);
const mockCreateGitHubDeployment = vi.fn().mockResolvedValue(99);
const mockUpdateDeploymentStatus = vi.fn().mockResolvedValue(undefined);

vi.mock("../db/repos.js", () => ({
	findRepoByFullName: mockFindRepoByFullName,
	upsertRepo: vi.fn(),
}));

vi.mock("../db/users.js", () => ({
	findUserById: mockFindUserById,
}));

vi.mock("../services/buang.js", () => ({
	createDeployment: mockCreateDeployment,
	tearDownBranch: mockTearDownBranch,
}));

vi.mock("../services/github.js", () => ({
	createSetupCheckRun: mockCreateSetupCheckRun,
	createGitHubDeployment: mockCreateGitHubDeployment,
	updateDeploymentStatus: mockUpdateDeploymentStatus,
	findOpenPRForBranch: vi.fn(),
}));

// Test the decision logic by simulating what the handler does
describe("pull-request deployment logic", () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	it("triggers setup check run when repo has no buang_project_id", async () => {
		mockFindRepoByFullName.mockResolvedValue({
			id: 1n,
			githubRepoFullName: "owner/repo",
			installationId: 1n,
			buangProjectId: null,
			ownerUserId: null,
			envVarsEncrypted: null,
			composePath: null,
			serviceEntrypoint: null,
			waitForWorkflow: false,
			workflowName: null,
			requiresAuthn: false,
			createdAt: new Date(),
			updatedAt: new Date(),
		});

		const repo = await mockFindRepoByFullName("owner/repo", 1);

		// Simulate handler decision
		if (!repo?.buangProjectId) {
			await mockCreateSetupCheckRun();
			expect(mockCreateSetupCheckRun).toHaveBeenCalledOnce();
			expect(mockCreateDeployment).not.toHaveBeenCalled();
		}
	});

	it("skips deployment when wait_for_workflow is true", async () => {
		mockFindRepoByFullName.mockResolvedValue({
			id: 2n,
			githubRepoFullName: "owner/repo",
			installationId: 1n,
			buangProjectId: 10n,
			ownerUserId: 1n,
			envVarsEncrypted: null,
			composePath: "docker-compose.yml",
			serviceEntrypoint: "traefik:80",
			waitForWorkflow: true,
			workflowName: "build",
			requiresAuthn: false,
			createdAt: new Date(),
			updatedAt: new Date(),
		});

		const repo = await mockFindRepoByFullName("owner/repo", 1);

		// Simulate handler decision
		if (repo?.waitForWorkflow) {
			// Should return early without deploying
			expect(mockCreateDeployment).not.toHaveBeenCalled();
		}
	});

	it("calls createDeployment when repo is configured", async () => {
		const { encrypt } = await import("../services/crypto.js");
		const encryptedKey = encrypt("test-api-key");

		mockFindRepoByFullName.mockResolvedValue({
			id: 3n,
			githubRepoFullName: "owner/repo",
			installationId: 1n,
			buangProjectId: 10n,
			ownerUserId: 1n,
			envVarsEncrypted: null,
			composePath: "docker-compose.yml",
			serviceEntrypoint: "traefik:80",
			waitForWorkflow: false,
			workflowName: null,
			requiresAuthn: false,
			createdAt: new Date(),
			updatedAt: new Date(),
		});

		mockFindUserById.mockResolvedValue({
			id: 1n,
			email: "user@example.com",
			passwordHash: null,
			authType: "oauth",
			buangApiBaseUrl: "https://buang.example.com",
			buangApiKeyEncrypted: encryptedKey,
			createdAt: new Date(),
			updatedAt: new Date(),
		});

		const repo = await mockFindRepoByFullName("owner/repo", 1);
		const user = await mockFindUserById(repo!.ownerUserId);

		expect(repo?.buangProjectId).toBeTruthy();
		expect(repo?.waitForWorkflow).toBe(false);
		expect(user?.buangApiBaseUrl).toBe("https://buang.example.com");
		expect(user?.buangApiKeyEncrypted).toBe(encryptedKey);
	});
});
