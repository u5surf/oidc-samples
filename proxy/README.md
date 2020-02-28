# MIRACL Proxy

Proxy implementation that allows us to modify the traffic for testing purposes.

`HOST` - Host to listen on. The default is "0.0.0.0".

`PORT` - Port of the listening host. The default is "8080".

`VERBOSE`- Log all requests to stdout. The default value is false.

In order to use it locally, you need to execute the following commands:

```
docker build -t proxy .
docker run -p 8080:8080 proxy
```
