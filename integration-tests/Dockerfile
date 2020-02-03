FROM docker.pkg.github.com/miracl/oidc-samples/test-runner:latest
WORKDIR /src
COPY . .

RUN mkdir -p ./vendor/code.miracl.com/mfa/pkg
RUN mv /gomiracl ./vendor/code.miracl.com/mfa/pkg
RUN go test -mod=vendor -c -o integration-tests .
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ./entrypoint.sh
