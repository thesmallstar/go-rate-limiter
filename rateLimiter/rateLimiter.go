package rateLimiter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/time/rate"
)

const DEFAULT_LIMIT_KEY = "DEFAULT ANY"

type RateLimiter struct {
	limiters map[string]map[string]*rate.Limiter
}

func GetRateLimiter() *RateLimiter {
	newRateLimiter := new(RateLimiter)
	newRateLimiter.limiters = make(map[string]map[string]*rate.Limiter)
	return newRateLimiter
}

func (r *RateLimiter) LoadConfig(path string) {

	var config ConfigFile

	file, err := os.Open(path)
	fmt.Println((err))
	defer file.Close()

	byteConfig, _ := ioutil.ReadAll(file)
	_ = json.Unmarshal(byteConfig, &config)

	for _, user := range config.Users {
		r.limiters[user.UserID] = make(map[string]*rate.Limiter)
		for _, route := range user.Limits {
			r.limiters[user.UserID][route.Route+" "+route.TypeReq] = rate.NewLimiter(rate.Limit(route.Rate), route.Limit)
		}
	}

	fmt.Println((r.limiters))
}

func (r *RateLimiter) IsValidRequest(userID string, route string, reqType string) bool {
	//	print(userID)
	if user, ok := r.limiters[userID]; ok {
		//	print("here-1")
		if limiter, ok := user[route+" "+reqType]; ok {
			//	print("here-2")
			return limiter.Allow()
		} else {
			//	print("here-3")
			if limiter, ok := user[DEFAULT_LIMIT_KEY]; ok {
				return limiter.Allow()
			} else {
				return false
			}
		}
	}
	return false
}
