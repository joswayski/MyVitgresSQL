package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func main() {
	targetURL, err := url.Parse("https://api.averagedatabase.com")
	if err != nil {
		log.Fatal("yolo")
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		shard := rand.Intn(20001) - 10000

		req.Host = targetURL.Host
		req.URL.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme

		req.Header.Set("shard", strconv.Itoa(shard))

		if apiKey := req.Header.Get("x-averagedb-api-key"); apiKey != "" {
			req.Header.Set("x-averagedb-api-key", apiKey)
		}

		if req.Header.Get("User-Agent") == "" {
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; API-Proxy/1.0)")
		}

		req.Header.Del("X-Forwarded-For")
		req.Header.Del("X-Forwarded-Proto")
		req.Header.Del("X-Real-IP")

	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Server")
		return nil
	}

	server := &http.Server{
		Addr:    ":6969",
		Handler: proxy,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
