FROM docker.pkg.github.com/miracl/oidc-samples/proxy:latest as proxy

FROM node:alpine
ENV CA_CERT /usr/local/share/ca-certificates/ca.pem
COPY --from=proxy /src/ca.crt &CA_CERT
COPY . /src

RUN cd /src && npm install

ENTRYPOINT ["node", "index.js"]
