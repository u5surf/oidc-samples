FROM golang:latest

COPY . /src
RUN cd /src \
    && go build -o /bin/proxy . \
    && cp entrypoint.sh /bin/

ENTRYPOINT entrypoint.sh
