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

The proxy has two modes - one that simply forwards the traffic and one which allows us to modify the requests and responses coming ang going through the proxy, using Man-In-The-Middle (MITM).

The way you can use the MITM proxy is by making a `POST` request to the `/session` endpoint. The JSON payload of that request must contain a `mediatorUrl` property, which is an endpoint in another system (in this case, the integration tests), to which the proxy will redirect all traffic and the response of which will contain the modified payload.

The goal here is for an external system to be able to modify the traffic between two communicating systems. 

There is also a way to stop the MITM session and you do this, by making an empty `DELETE` request to the `/session` endpoint. When you do this, it will continue proxying, but just by forwarding the requests/responses.
