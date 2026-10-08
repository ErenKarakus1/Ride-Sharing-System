package proxy

import (
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

	return &Proxy{
		proxy: httputil.NewSingleHostReverseProxy(targetURL),
	}
}

func (p *Proxy) Handler() http.Handler {
	return p.proxy
}
