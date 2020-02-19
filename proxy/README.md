# MIRACL Proxy

We have a proxy implementation using the [goproxy](https://github.com/elazarl/goproxy) library. We need this in order to be able to tamper with requests and responses made between the samples and our platform. We also use it to test if the OIDC libraries which we use in all samples work through a proxy.

`HOST` - Host to listen on. The default is "0.0.0.0".

`PORT` - Port of the listening host. The default is "8080".

`VERBOSE`- Log all requests to stdout. The default value is false.

In order to use it locally, you need to execute the following commands:

```
docker build -t proxy .
docker run -p 8080:8080 proxy
```
