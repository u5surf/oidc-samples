# MIRACL Trust OIDC Integration Samples

![go](https://github.com/miracl/oidc-samples/workflows/go/badge.svg)
![nodejs](https://github.com/miracl/oidc-samples/workflows/nodejs/badge.svg)


These samples provide example integrations between the [MIRACL Trust](https://miracl.com) platform and various OIDC libraries.

## ENV Variables

All samples work with the following environment variables:

`HOST` - Host to listen on. The default is "localhost".

`PORT` - Port of the listening host. The default is "8000".

`ISSUER` - OpenID Connect Issuer. The default is "https://api.mpin.io".

`REDIRECT_URL` - The redirect URL of the app in the MIRACL Trust platform. The default value is "http://localhost:8000/login".

`CLIENT_ID` - The client id of the app in the MIRACL Trust platform. Has no default value and is required.

`CLIENT_SECRET`- The client secret of the app in the MIRACL Trust platform. Has no default value and is required.

`PROXY_HOST`- The host address of the proxy, which you wish to run the sample behind. The default value is empty string.

`PROXY_PORT`- The port of the proxy, which you wish to run the sample behind. The default value is empty string.


The required env vars are `CLIENT_ID` and `CLIENT_SECRET`.

To get those values, you'll need to [register and create an app in our platform](https://docs.miracl.cloud/get-started/).

We recommend you to leave the `ISSUER` to the default value. This way you'll be using [https://trust.miracl.cloud](https://trust.miracl.cloud), the full docs of which you can find [here](https://docs.miracl.cloud/).

## Usage

You can run any sample as Docker container

```
cd samples/<variant>
docker build -t sample .
docker run -p 8000:8000 -e CLIENT_ID=<client-id> -e CLIENT_SECRET=<client-secret> sample
```

Afterwards visit `http://localhost:8000`. If you need to set another `PORT` or `HOST`, change the `8000:8000` port mapping and add the env variables in the `docker run` command.

## Running through proxy

In order to test how OIDC libraries behave in some edge cases(for ex. when the OIDC server misbehaves) - we need to modify the traffic between the library and the sample showcasing that library.

You have the option to use our proxy with the provided samples. You can check the [README](proxy/README.md) in the proxy directory for how to build and run it.

Provided that you have a built docker image of the proxy and the sample that you wish to run, it's as simple as just running both `docker run` commands, with the addition of the `PROXY_HOST` and `PROXY_PORT` environment variables. If you've stuck to the default values, those commands are:

```
docker run -p 8080:8080 proxy
docker run -p 8000:8000 -e PROXY_HOST=127.0.0.1 -e PROXY_PORT=8080 -e CLIENT_ID=<client-id> -e CLIENT_SECRET=<client-secret> sample
```

You can confirm that the requests from the sample are going through a proxy if you enable the verbose mode of the proxy, using the `VERBOSE` environment variable in the command above. When the proxy and sample are started and you complete a registration and authentication, the proxy output will log out the information of the proxied requests.
