DROP TRIGGER IF EXISTS repositories_set_updated_at ON repositories;
DROP INDEX IF EXISTS idx_repos_full_name_installation;
DROP TABLE IF EXISTS repositories;
