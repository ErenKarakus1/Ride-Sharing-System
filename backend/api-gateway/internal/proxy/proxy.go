package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	proxy *httputil.ReverseProxy
}

func New(target string) *Proxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		log.Printf("proxy error target=%s path=%s request_id=%s error=%v", targetURL.String(), request.URL.Path, request.Header.Get("X-Request-ID"), err)
		http.Error(writer, "upstream service unavailable", http.StatusBadGateway)
	}

	return &Proxy{proxy: reverseProxy}
}

func (p *Proxy) Handler() http.Handler {
	return p.proxy
}
