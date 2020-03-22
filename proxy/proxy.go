package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/elazarl/goproxy"
)

type proxy struct {
	sessionMux sync.RWMutex
	session    *session
}

type session struct {
	ModifyURL string `json:"modifyUrl"`
}

func (p *proxy) setupProxy() *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()

	proxy.NonproxyHandler = p.sessionHandler()

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	proxy.OnResponse().DoFunc(p.modifyResponse)

	return proxy
}

func (p *proxy) sessionHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/session" {
			http.Error(w, "Error: Only requests to /session allowed", http.StatusNotFound)
			return
		}
		if req.Method != http.MethodPost && req.Method != http.MethodDelete {
			http.Error(w, "Error: Only POST and DELETE allowed", http.StatusMethodNotAllowed)
			return
		}
		if req.Method == http.MethodPost && p.session != nil {
			http.Error(w, "Error: There is already a started session", http.StatusBadRequest)
			return
		}
		if req.Method == http.MethodDelete && p.session == nil {
			http.Error(w, "Error: There is no started session", http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodPost {
			p.startSession(w, req)
		}
		if req.Method == http.MethodDelete {
			p.sessionMux.Lock()
			p.session = nil
			p.sessionMux.Unlock()
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (p *proxy) startSession(w http.ResponseWriter, req *http.Request) error {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	defer req.Body.Close()

	var sessionResp session
	if err = json.Unmarshal(reqBody, &sessionResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	p.sessionMux.Lock()
	p.session = &session{
		ModifyURL: sessionResp.ModifyURL,
	}
	p.sessionMux.Unlock()

	w.WriteHeader(http.StatusOK)

	return nil
}

func (p *proxy) modifyResponse(proxiedResp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if p.session == nil {
		return proxiedResp
	}

	client := &http.Client{}

	p.sessionMux.RLock()
	modifyReq, _ := http.NewRequest(proxiedResp.Request.Method, p.session.ModifyURL, nil)
	p.sessionMux.RUnlock()

	// copy body and headers from the response that the proxy intercepted
	// to the place where they can be modified and returned back as response
	modifyReq.Header = proxiedResp.Header
	modifyReq.Body = proxiedResp.Body
	modifiedResp, err := client.Do(modifyReq)

	if err != nil {
		log.Fatal("Error sending request to modifier", err)
	}

	proxiedResp.Body = ioutil.NopCloser(modifiedResp.Body)
	proxiedResp.Header = modifiedResp.Header

	return proxiedResp
}
