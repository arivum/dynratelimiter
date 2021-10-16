FROM golang:1.17.1-alpine3.14 AS builder
RUN apk update && \
    apk add make git clang linux-headers musl-dev
COPY . /src
WORKDIR /src
RUN make build-dynratelimiter-operator

FROM scratch
COPY --from=builder /src/dynratelimiter-operator /usr/bin/dynratelimiter-operator
ENTRYPOINT [ "/usr/bin/dynratelimiter-operator" ]