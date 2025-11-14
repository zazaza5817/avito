import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';

// === КОНФИГУРАЦИЯ ===
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const ADMIN_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYWRtaW4tdXNlci1pZCIsImlzX2FkbWluIjp0cnVlLCJleHAiOjE3OTQ2ODI4MTMsImlhdCI6MTc2MzE0NjgxM30.juINEX89nV19P0aY4gi2W8MgLBw8lsbX38n-BnccLDA';

// === МЕТРИКИ ДЛЯ SLI ===
const errorRate = new Rate('sli_error_rate');
const availability = new Rate('sli_availability');
const latencyP95 = new Trend('sli_latency_p95');
const latencyP99 = new Trend('sli_latency_p99');

// Метрики по эндпоинтам
const endpointMetrics = {
  health: { latency: new Trend('health_latency'), errors: new Counter('health_errors') },
  teamGet: { latency: new Trend('team_get_latency'), errors: new Counter('team_get_errors') },
  teamAdd: { latency: new Trend('team_add_latency'), errors: new Counter('team_add_errors') },
  userGetReview: { latency: new Trend('user_getreview_latency'), errors: new Counter('user_getreview_errors') },
  userSetActive: { latency: new Trend('user_setactive_latency'), errors: new Counter('user_setactive_errors') },
  prCreate: { latency: new Trend('pr_create_latency'), errors: new Counter('pr_create_errors') },
  prMerge: { latency: new Trend('pr_merge_latency'), errors: new Counter('pr_merge_errors') },
  prReassign: { latency: new Trend('pr_reassign_latency'), errors: new Counter('pr_reassign_errors') },
};

// === НАСТРОЙКИ ТЕСТА ===
export const options = {
  stages: [
    { duration: '5s', target: 100 },
    { duration: '1m', target: 100 },
    { duration: '5s', target: 100 }
  ],
  thresholds: {
    // SLI: Availability >= 99.9%
    'sli_availability': ['rate>0.999'],
    // SLI: P95 Latency < 300ms
    'sli_latency_p95': ['p(95)<300'],
    // SLI: P99 Latency < 500ms
    'sli_latency_p99': ['p(99)<500'],
    // SLI: Error Rate < 0.1%
    'sli_error_rate': ['rate<0.001'],
  },
};

function makeRequest(method, url, body, token, metric, allowedErrorCodes = []) {
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  if (token) {
    params.headers['Authorization'] = `Bearer ${token}`;
  }

  const res = method === 'GET' 
    ? http.get(url, params)
    : http.post(url, body, params);

  metric.latency.add(res.timings.duration);
  latencyP95.add(res.timings.duration);
  latencyP99.add(res.timings.duration);

  // Успех: 2xx-3xx или ожидаемый код ошибки (например 409)
  const isAllowedError = allowedErrorCodes.includes(res.status);
  const success = (res.status >= 200 && res.status < 400) || isAllowedError;
  
  if (!success) {
    metric.errors.add(1);
    errorRate.add(1);
    availability.add(0);
  } else {
    errorRate.add(0);
    availability.add(1);
  }

  return { res, success };
}

// === ГЛАВНЫЙ СЦЕНАРИЙ ===
export default function () {
  const vu = __VU;
  const iter = __ITER;
  const timestamp = Date.now();
  const uniqueId = `${vu}-${iter}-${timestamp}`;

  // 1. Health Check
  makeRequest('GET', `${BASE_URL}/health`, null, null, endpointMetrics.health);
  sleep(0.1);

  // 2. Получение команды
  makeRequest('GET', `${BASE_URL}/team/get?team_name=test_backend`, null, ADMIN_TOKEN, endpointMetrics.teamGet);
  sleep(0.1);

  // 3. Получение PR пользователя
  const testUserId = `tb${(vu % 18) + 1}`;
  makeRequest('GET', `${BASE_URL}/users/getReview?user_id=${testUserId}`, null, ADMIN_TOKEN, endpointMetrics.userGetReview);
  sleep(0.1);

  // 4. Создание команды
  const teamPayload = JSON.stringify({
    team_name: `load-test-team-${uniqueId}`,
    members: [
      { user_id: `lt-${uniqueId}-1`, username: `LoadTestUser1-${vu}`, is_active: true },
      { user_id: `lt-${uniqueId}-2`, username: `LoadTestUser2-${vu}`, is_active: true },
      { user_id: `lt-${uniqueId}-3`, username: `LoadTestUser3-${vu}`, is_active: true },
    ],
  });
  makeRequest('POST', `${BASE_URL}/team/add`, teamPayload, ADMIN_TOKEN, endpointMetrics.teamAdd);
  sleep(0.2);

  // 5. Создание PR (admin)
  const prId = `pr-load-${uniqueId}`;
  const authorId = `tb${(vu % 18) + 1}`; // tb1-tb18 (активные)
  const prPayload = JSON.stringify({
    pull_request_id: prId,
    pull_request_name: `Load Test PR ${uniqueId}`,
    author_id: authorId,
  });
  const prResult = makeRequest('POST', `${BASE_URL}/pullRequest/create`, prPayload, ADMIN_TOKEN, endpointMetrics.prCreate);
  sleep(0.2);

  if (prResult.success && prResult.res.status === 201) {
    try {
      const prData = JSON.parse(prResult.res.body);
      const assignedReviewers = prData.pr?.assigned_reviewers || [];
      
      // 6. Переназначение ревьювера (409 ожидается если нет кандидатов)
      if (assignedReviewers.length > 0) {
        const oldReviewerId = assignedReviewers[0];
        const reassignPayload = JSON.stringify({
          pull_request_id: prId,
          old_user_id: oldReviewerId,
        });
        const reassignResult = makeRequest('POST', `${BASE_URL}/pullRequest/reassign`, reassignPayload, ADMIN_TOKEN, endpointMetrics.prReassign, [409]);
        sleep(0.2);
        
        if (reassignResult.success && reassignResult.res.status === 200) {
          try {
            const reassignData = JSON.parse(reassignResult.res.body);
            const newReviewerId = reassignData.replaced_by;
          } catch (e) {
          }
        }
      }
    } catch (e) {
    }

    const mergePayload = JSON.stringify({
      pull_request_id: prId,
    });
    makeRequest('POST', `${BASE_URL}/pullRequest/merge`, mergePayload, ADMIN_TOKEN, endpointMetrics.prMerge);
    sleep(0.2);
  }

  // 7. Изменение статуса пользователя 
  const userToToggle = `tb${(vu % 18) + 1}`;
  const activePayload = JSON.stringify({
    user_id: userToToggle,
    is_active: Math.random() > 0.5,
  });
  makeRequest('POST', `${BASE_URL}/users/setIsActive`, activePayload, ADMIN_TOKEN, endpointMetrics.userSetActive);
  sleep(0.2);
}
