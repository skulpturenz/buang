CREATE TABLE IF NOT EXISTS repositories (
    id                    BIGSERIAL PRIMARY KEY,
    github_repo_full_name TEXT NOT NULL,
    installation_id       BIGINT NOT NULL,
    buang_project_id      BIGINT,
    owner_user_id         BIGINT REFERENCES users(id) ON DELETE SET NULL,
    env_vars_encrypted    TEXT,
    compose_path          TEXT,
    service_entrypoint    TEXT,
    wait_for_workflow     BOOLEAN NOT NULL DEFAULT FALSE,
    workflow_name         TEXT,
    requires_authn        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_repos_full_name_installation
    ON repositories (github_repo_full_name, installation_id);

CREATE TRIGGER repositories_set_updated_at
    BEFORE UPDATE ON repositories
    FOR EACH ROW
    EXECUTE PROCEDURE set_updated_at();
