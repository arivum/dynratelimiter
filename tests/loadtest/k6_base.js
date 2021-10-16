/*
 * Copyright (c) 2021, arivum.
 * All rights reserved.
 * SPDX-License-Identifier: MIT
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/MIT
 */

import http from "k6/http";
import { check, group, sleep } from "k6";

export let options = {
    stages: [
        { duration: "3m", target: 100 }, // simulate ramp-up of traffic from 1 to 100 users over 5 minutes.
        { duration: "14m", target: 100 }, // stay at 100 users for 10 minutes
        { duration: "3m", target: 0 }, // ramp-down to 0 users
    ],
};

const BASE_URL = `http://${__ENV.TARGET_IP}:8080`;

export default () => {
    let response = http.get(`${BASE_URL}/`);

    check(response, {
        "sucessfully got answer": (resp) => resp.body && resp.body.startsWith("OK"),
    });
    sleep(1);
};