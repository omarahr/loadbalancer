package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	targetServers = []string{
		"http://localhost:8080/",
		"http://localhost:8081/",
		"http://localhost:8082/",
	}
)

func getProxies() []*Proxy {
	var proxies []*Proxy

	for _, targetS := range targetServers {
		target, err := url.Parse(targetS)
		if err != nil {
			log.Fatal(err)
		}

		proxies = append(proxies, &Proxy{
			Target: target,
			Proxy:  httputil.NewSingleHostReverseProxy(target),
			client: &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(target),
				},
				Timeout: 1 * time.Second,
			},
			healthy: true,
		})
	}

	return proxies
}

type ProxyList struct {
	mu           sync.Mutex
	currentProxy int32
	proxies      []*Proxy
}

type Proxy struct {
	Target  *url.URL
	Proxy   *httputil.ReverseProxy
	client  *http.Client
	healthy bool
}

func (p *Proxy) Healthy() bool {
	resp, err := p.client.Get(p.Target.String() + "health")
	if err != nil {
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func NewProxyList() *ProxyList {
	return &ProxyList{
		currentProxy: 0,
		proxies:      getProxies(),
	}
}

func (pl *ProxyList) Next() *Proxy {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.currentProxy++
	if int(pl.currentProxy) >= len(pl.proxies) {
		pl.currentProxy = 0
	}

	return pl.proxies[pl.currentProxy]
}

func (pl *ProxyList) BigNext() *Proxy {
	var next *Proxy

	for next == nil || !next.healthy {
		next = pl.Next()
	}

	return next
}

func (pl *ProxyList) UpdateHealth(index int, healthy bool) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.proxies[index].healthy = healthy
}

func HealthChecker(proxyList *ProxyList) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped properly

	// Infinite loop to keep the program running
	for {
		select {
		case <-ticker.C:
			for idx, proxy := range proxyList.proxies {
				proxyList.UpdateHealth(idx, proxy.Healthy())
			}
		}
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	proxyList := NewProxyList()

	go HealthChecker(proxyList)

	var reqCounter int32

	r.Any("/*proxyPath", func(c *gin.Context) {
		// Update the request URL to the target URL
		nextProxy := proxyList.BigNext()
		nextCounter := atomic.AddInt32(&reqCounter, 1)

		c.Request.URL.Scheme = nextProxy.Target.Scheme
		c.Request.URL.Host = nextProxy.Target.Host
		c.Request.URL.Path = c.Param("proxyPath")

		additionalQuery := "reqCounter=" + strconv.Itoa(int(nextCounter))

		// Append the additional query parameter
		originalQuery := c.Request.URL.RawQuery
		if originalQuery == "" {
			c.Request.URL.RawQuery = additionalQuery
		} else {
			c.Request.URL.RawQuery = originalQuery + "&" + additionalQuery
		}

		// Use the reverse proxy to handle the request
		nextProxy.Proxy.ServeHTTP(c.Writer, c.Request)
	})

	if err := r.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
