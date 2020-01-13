# MIRACL Trust OIDC Integration Samples

![](https://github.com/miracl/oidc-samples/workflows/Go-sample/badge.svg)

## OIDC Credentials

To get client ID and secret check our documentation [here](https://docs.miracl.cloud/get-started/#client-id-and-secret).

## Run a sample as Docker container

```
cd samples/<variant>
docker build -t sample .
docker run -p 8000:8000 sample -client-id <client-id> -client-secret <client-secret>
```
