FROM docker.pkg.github.com/miracl/oidc-samples/test-runner:latest
WORKDIR /src
COPY . .

RUN mkdir -p ./vendor/code.miracl.com/mfa/pkg \
    && mv /gomiracl ./vendor/code.miracl.com/mfa/pkg \
    && go test -mod=vendor -c -o integration-tests . \
    && chmod +x ./entrypoint.sh

ENTRYPOINT ./entrypoint.sh
