FROM golang:alpine
WORKDIR /src
COPY . .

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY ./ca.crt /usr/local/share/ca-certificates/ca.pem

RUN update-ca-certificates

RUN go build -o sample .
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ./entrypoint.sh
