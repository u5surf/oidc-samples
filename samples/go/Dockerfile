FROM docker.pkg.github.com/miracl/oidc-samples/proxy:latest as proxy


FROM golang:alpine as builder
COPY . /src
RUN cd /src \
    && go build -o /bin/sample .


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=proxy /src/ca.crt /usr/local/share/ca-certificates/ca.pem
COPY --from=builder /bin/sample /bin/sample
COPY --from=builder /src/entrypoint.sh /bin/entrypoint.sh

ENTRYPOINT entrypoint.sh
