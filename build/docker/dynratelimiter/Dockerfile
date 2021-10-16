FROM golang:1.17.1-alpine3.14 AS builder
RUN apk update && \
    apk add make git clang linux-headers musl-dev
COPY . /src
WORKDIR /src
RUN make

FROM alpine:3.14
COPY --from=builder /src/dynratelimiter /usr/bin/dynratelimiter
ENTRYPOINT [ "/usr/bin/dynratelimiter" ]