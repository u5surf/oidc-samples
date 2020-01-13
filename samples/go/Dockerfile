FROM golang:alpine
WORKDIR /src
COPY . .
RUN go build -o sample .
ENTRYPOINT ["./sample"]
