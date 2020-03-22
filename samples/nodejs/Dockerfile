FROM docker.pkg.github.com/miracl/oidc-samples/proxy:latest as proxy


FROM node:alpine

COPY . /src

ENV CA_CERT /src/ca.crt
COPY --from=proxy /src/ca.crt $CA_CERT

RUN cd /src && npm install

ENTRYPOINT ["node", "src/index.js"]
