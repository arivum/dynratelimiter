# Tested Benchmarks

## Build

```bash
make build build-stressme build-top
```

## Run

* Start benchmark environment
    ```bash
    docker run --privileged -it --name bench --rm -v $PWD:/export --workdir=/export --cpus="0.25" -m 5G golang:1.17.1-alpine3.14
    ```
* Run resource-top to control resource allocation
    ```bash
    docker exec -it bench /export/resource-top
    ```
* Start stressme server
    ```bash
    docker exec -it bench /export/stressme -cpu
    ```
### Baseline without rate limiting
* Bomb with curl
    ```bash
    OUT_FILE=count_without_bpf
    CONTAINER_IP=$(docker inspect bench --format='{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' | head -n 1)
    for i in {0..5}; do for j in {0..100}; do (curl --connect-timeout 1 ${CONTAINER_IP}:8080 >/dev/null 2>&1; echo $? >> ${OUT_FILE} ) &; done; sleep 3; done;

    # Once finished
    echo $(cat ${OUT_FILE} | grep ^0$ | wc -l) / $(cat ${OUT_FILE} | wc -l)
    # e.g. 36 / 606
    ```

### With rate limiting
* Start dynamic rate limiter
    ```bash
    docker exec -it bench /export/dynratelimiter -conf /export/examples/dynratelimit.yaml
    ```
* Bomb with curl
    ```bash
    OUT_FILE=count_bpf
    CONTAINER_IP=$(docker inspect bench --format='{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' | head -n 1)
    for i in {0..5}; do for j in {0..100}; do (curl --connect-timeout 1 ${CONTAINER_IP}:8080 >/dev/null 2>&1; echo $? >> ${OUT_FILE} ) &; done; sleep 3; done;

    # Once finished
    echo $(cat ${OUT_FILE} | grep ^0$ | wc -l) / $(cat ${OUT_FILE} | wc -l)
    # e.g. 171 / 606
    ```