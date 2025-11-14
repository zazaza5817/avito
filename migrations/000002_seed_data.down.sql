-- Откат seed

-- Удаляем тестовые PR и их ревьюверов
DELETE FROM pr_reviewers WHERE pull_request_id IN (
  'pr-test-1', 'pr-test-2', 'pr-test-3', 'pr-test-4', 'pr-test-5', 'pr-test-6'
);
DELETE FROM pull_requests WHERE pull_request_id IN (
  'pr-test-1', 'pr-test-2', 'pr-test-3', 'pr-test-4', 'pr-test-5', 'pr-test-6'
);

-- Удаляем пользователей test_backend (tb1-tb20)
DELETE FROM users WHERE user_id IN (
  'tb1', 'tb2', 'tb3', 'tb4', 'tb5', 'tb6', 'tb7', 'tb8', 'tb9', 'tb10',
  'tb11', 'tb12', 'tb13', 'tb14', 'tb15', 'tb16', 'tb17', 'tb18', 'tb19', 'tb20'
);

-- Удаляем пользователей test_frontend (tf1-tf20)
DELETE FROM users WHERE user_id IN (
  'tf1', 'tf2', 'tf3', 'tf4', 'tf5', 'tf6', 'tf7', 'tf8', 'tf9', 'tf10',
  'tf11', 'tf12', 'tf13', 'tf14', 'tf15', 'tf16', 'tf17', 'tf18', 'tf19', 'tf20'
);

-- Удаляем команды
DELETE FROM teams WHERE team_name IN ('test_backend', 'test_frontend');
