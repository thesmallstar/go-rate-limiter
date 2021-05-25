package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thesmallstar/go-rate-limiter/go-rate-limiter/rateLimiter"
)

type RateLimitMiddleWare struct {
	handler           http.Handler
	rateLimiterObject *rateLimiter.RateLimiter
}

func (l *RateLimitMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL.Path + " " + r.Method + " BY " + r.Header.Get("USERID"))
	userID := r.Header.Get("USERID")
	route := r.URL.Path
	reqType := r.Method

	if !(l.rateLimiterObject.IsValidRequest(userID, route, reqType)) {
		http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
		return
	}
	l.handler.ServeHTTP(w, r)

}

func newRateLimitMiddleWare(handlerToWrap http.Handler) *RateLimitMiddleWare {
	r := rateLimiter.GetRateLimiter()
	r.LoadConfig("rateLimiter/config.json")
	return &RateLimitMiddleWare{handlerToWrap, r}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func CurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now().Format(time.Kitchen)
	w.Write([]byte(fmt.Sprintf("the current time is %v", curTime)))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/hello", HelloHandler)
	mux.HandleFunc("/v1/time", CurrentTimeHandler)

	wrappedMux := newRateLimitMiddleWare(mux)

	log.Printf("server is listening at %s", "localhost:4000")

	log.Fatal(http.ListenAndServe(":4000", wrappedMux))

}
