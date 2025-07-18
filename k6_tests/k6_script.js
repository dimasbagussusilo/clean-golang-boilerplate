import http from 'k6/http';
import {check, group, sleep} from 'k6';

export const options = {
    stages: [
        {duration: '30s', target: 150},
        {duration: '1m', target: 150},
        {duration: '10s', target: 0},
    ],
    thresholds: {

        'http_req_failed': ['rate<0.02'],
        'http_req_duration{api:gin}': ['p(95)<600'],
        'http_req_duration{api:fiber}': ['p(95)<600'],
    },
};

const payload = JSON.stringify({
    name: 'testing store',
    description: 'only for performance testing',
    type: 'retail',
});

const params = {
    headers: {
        'Content-Type': 'application/json',
    },
};

export default function () {
    const GIN_URL = 'http://localhost:8080/v1/stores';
    const FIBER_URL = 'http://localhost:9090/v1/stores';

    group('Gin POST /v1/stores', function () {
        const res = http.post(GIN_URL, payload, {...params, tags: {api: 'gin'}});
        check(res, {
            'status is 201': (r) => r.status === 201,
        });
    });

    group('Fiber (Prefork) POST /v1/stores', function () {
        const res = http.post(FIBER_URL, payload, {...params, tags: {api: 'fiber'}});
        check(res, {
            'status is 201': (r) => r.status === 201,
        });
    });

    sleep(1);
}