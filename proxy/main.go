package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	setCA(caCert, caKey)

	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.Logf("%v", "We can see what APIs are being called!")
		return req, ctx.Resp
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		ctx.Logf("%v", "We can modify some data coming back!")
		return resp
	})

	var host, port string
	flag.StringVar(&host, "host", "0.0.0.0", "Listen host")
	flag.StringVar(&port, "port", "8080", "Listen port")
	verbose := flag.Bool("v", false, "Log every proxied request to stdout")
	flag.Parse()

	proxy.Verbose = *verbose
	address := host + ":" + port

	log.Fatal(http.ListenAndServe(address, proxy))
}
