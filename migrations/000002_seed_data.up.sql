-- Seed данных: команды, пользователи, pull_requests и назначения ревьюверов

-- ========================================
-- ТЕСТОВЫЕ КОМАНДЫ
-- ========================================
INSERT INTO teams (team_name) VALUES ('test_backend') ON CONFLICT (team_name) DO NOTHING;
INSERT INTO teams (team_name) VALUES ('test_frontend') ON CONFLICT (team_name) DO NOTHING;

-- ========================================
-- КОМАНДА test_backend - 20 человек
-- ========================================
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb1', 'BackendUser1', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb2', 'BackendUser2', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb3', 'BackendUser3', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb4', 'BackendUser4', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb5', 'BackendUser5', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb6', 'BackendUser6', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb7', 'BackendUser7', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb8', 'BackendUser8', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb9', 'BackendUser9', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb10', 'BackendUser10', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb11', 'BackendUser11', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb12', 'BackendUser12', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb13', 'BackendUser13', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb14', 'BackendUser14', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb15', 'BackendUser15', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb16', 'BackendUser16', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb17', 'BackendUser17', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb18', 'BackendUser18', 'test_backend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb19', 'BackendUser19', 'test_backend', false) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tb20', 'BackendUser20', 'test_backend', false) ON CONFLICT (user_id) DO NOTHING;

-- ========================================
-- КОМАНДА test_frontend - 20 человек
-- ========================================
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf1', 'FrontendUser1', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf2', 'FrontendUser2', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf3', 'FrontendUser3', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf4', 'FrontendUser4', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf5', 'FrontendUser5', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf6', 'FrontendUser6', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf7', 'FrontendUser7', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf8', 'FrontendUser8', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf9', 'FrontendUser9', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf10', 'FrontendUser10', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf11', 'FrontendUser11', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf12', 'FrontendUser12', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf13', 'FrontendUser13', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf14', 'FrontendUser14', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf15', 'FrontendUser15', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf16', 'FrontendUser16', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf17', 'FrontendUser17', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf18', 'FrontendUser18', 'test_frontend', true) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf19', 'FrontendUser19', 'test_frontend', false) ON CONFLICT (user_id) DO NOTHING;
INSERT INTO users (user_id, username, team_name, is_active) VALUES
  ('tf20', 'FrontendUser20', 'test_frontend', false) ON CONFLICT (user_id) DO NOTHING;

-- ========================================
-- ТЕСТОВЫЕ PULL REQUESTS для test_backend
-- ========================================
-- pr-test-1: OPEN, автор tb1, назначены tb2, tb3
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES ('pr-test-1', 'Add authentication', 'tb1', 'OPEN', NOW(), NULL)
ON CONFLICT (pull_request_id) DO NOTHING;

INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-1', 'tb2', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;
INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-1', 'tb3', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;

-- pr-test-2: OPEN, автор tb5, назначены tb6, tb7
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES ('pr-test-2', 'Fix database connection', 'tb5', 'OPEN', NOW(), NULL)
ON CONFLICT (pull_request_id) DO NOTHING;

INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-2', 'tb6', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;
INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-2', 'tb7', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;

-- pr-test-3: OPEN, автор tb10, назначены tb11, tb12
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES ('pr-test-3', 'Update dependencies', 'tb10', 'OPEN', NOW(), NULL)
ON CONFLICT (pull_request_id) DO NOTHING;

INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-3', 'tb11', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;
INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-3', 'tb12', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;

-- pr-test-4: MERGED, автор tb15
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES ('pr-test-4', 'Add logging', 'tb15', 'MERGED', NOW() - INTERVAL '1 day', NOW() - INTERVAL '6 hours')
ON CONFLICT (pull_request_id) DO NOTHING;

INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-4', 'tb16', NOW() - INTERVAL '1 day') ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;

-- ========================================
-- ТЕСТОВЫЕ PULL REQUESTS для test_frontend
-- ========================================
-- pr-test-5: OPEN, автор tf1, назначены tf2, tf3
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES ('pr-test-5', 'Redesign homepage', 'tf1', 'OPEN', NOW(), NULL)
ON CONFLICT (pull_request_id) DO NOTHING;

INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-5', 'tf2', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;
INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-5', 'tf3', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;

-- pr-test-6: OPEN, автор tf8, назначены tf9, tf10
INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
VALUES ('pr-test-6', 'Add mobile responsiveness', 'tf8', 'OPEN', NOW(), NULL)
ON CONFLICT (pull_request_id) DO NOTHING;

INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-6', 'tf9', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;
INSERT INTO pr_reviewers (pull_request_id, reviewer_id, assigned_at) VALUES
  ('pr-test-6', 'tf10', NOW()) ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING;

