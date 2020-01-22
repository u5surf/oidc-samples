FROM golang:alpine
WORKDIR /src
COPY . .

RUN go build -o sample .
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ./entrypoint.sh
