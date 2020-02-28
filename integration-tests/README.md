# MIRACL Trust Integration Tests for OIDC Samples

We have these integration tests in place, in order to test all our samples and the OIDC libraries' integration with the MIRACL Trust platform.

The integration tests use the following environment variables:

`API_URL` - The URL of the MPIN API. The default value is "http://api.mpin.io">

`SAMPLE_URL` - The URL of the sample that you want to test. The default value is "http://localhost:8000".

`REDIRECT_URL` - The redirect URL of the app in the MIRACL Trust platform. The default value is "http://localhost:8000/login".

`CLIENT_ID` - The client id of the app in the MIRACL Trust platform. Has no default value and is required.

`CLIENT_SECRET`- The client secret of the app in the MIRACL Trust platform. Has no default value and is required.

The required env vars are `CLIENT_ID` and `CLIENT_SECRET`.

To get those values, you'll need to [register and create an app in our platform](https://docs.miracl.cloud/get-started/).
__Note:__ Make sure to enable custom verification for the client that you're creating the app in. You can do this in the client settings section of the portal.

The only difference between running the samples and the tests, is that you should add `--network host` as a parameter. This meas that the container is binding its' ports directly to the host machine.

```
docker build -t int-tests .
docker run --network host -e CLIENT_ID=<client-id> -e CLIENT_SECRET=<client-secret> int-tests
```
