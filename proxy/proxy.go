package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	backendURL *url.URL
	proxy      *httputil.ReverseProxy
}

func NewProxy() *Proxy {
	backendURL, err := url.Parse("http://your-backend-service-url")
	if err != nil {
		panic(err)
	}
	return &Proxy{
		backendURL: backendURL,
		proxy:      httputil.NewSingleHostReverseProxy(backendURL),
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 这里可以进一步处理请求，如添加特定头信息
	r.Header.Set("X-Proxy-Info", "This is from proxy")
	p.proxy.ServeHTTP(w, r)
}
