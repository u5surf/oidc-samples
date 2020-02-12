FROM golang:latest
WORKDIR /src
COPY . .

RUN go build -o proxy .

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ./entrypoint.sh
