/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

import http from 'k6/http';

export let options = {
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 1000,
            timeUnit: '1s', // 1000 iterations per second, i.e. 1000 RPS
            duration: '10m',
            preAllocatedVUs: 100, // how large the initial pool of VUs would be
            maxVUs: 200, // if the preAllocatedVUs are not enough, we can initialize more
        },
    },
};

const BASE_URL = `http://${__ENV.TARGET_IP}:8080`;

export default function() {
    http.get(`${BASE_URL}/`);
}