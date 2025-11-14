-- Откат миграции: удаление таблиц в обратном порядке
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_team_name;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
