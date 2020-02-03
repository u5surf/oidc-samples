# MIRACL Trust OIDC Integration Samples

![go](https://github.com/miracl/oidc-samples/workflows/go/badge.svg)
![nodejs](https://github.com/miracl/oidc-samples/workflows/nodejs/badge.svg)


These samples provide example integrations between the [MIRACL Trust](https://miracl.com) platform and various OIDC libraries.

## ENV Variables

`HOST` - Host to listen on. The default is "localhost".
`PORT` - Port of the listening host. The default is "8000".
`ISSUER` - OpenID Connect Issuer. The default is "https://api.mpin.io"
`REDIRECT_URL` - The redirect URL of the app in the MIRACL Trust platform. The default value is "http://localhost:8000/login"
`CLIENT_ID` - The client id of the app in the MIRACL Trust platform. Has no default value and is required.
`CLIENT_SECRET`- The client secret of the app in the MIRACL Trust platform. Has no default value and is required.


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

Afterwards visit `http://localhost:8000`. If you need to set another port, change the `8000:8000` port mapping and add an `ADDR` env variable after the client ID and secret in `docker run`.
