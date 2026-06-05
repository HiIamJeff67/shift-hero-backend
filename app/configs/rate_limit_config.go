package configs

import (
	"time"

	"golang.org/x/time/rate"

	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

type RateLimitConfig struct {
	RateLimit         rate.Limit
	Burst             int
	UserLimit         int32
	WindowDuration    time.Duration
	BackendServerName types.BackendServerName
}

var DefaultAuthorizedRateLimitConfig = RateLimitConfig{
	RateLimit:         rate.Limit(300),                  // 300 requests/second
	Burst:             60,                               // allowed 60 additional requests/second for burst
	UserLimit:         9000,                             // 9000 requests/each life time of the bucket (= 9000 requests/`WindowDuration`) for each users
	WindowDuration:    time.Minute,                      // 1 minutes to reset the bucket
	BackendServerName: types.BackendServerName_EastAsia, // the current server
}

var DefaultUnauthorizedRateLimitConfig = RateLimitConfig{
	RateLimit:         rate.Limit(30),                   // 30 requests/second
	Burst:             15,                               // allowed 15 additional requests/second for burst
	UserLimit:         900,                              // 900 requests/each life time of the bucket (= 900 requests/`WindowDuration`) for each users
	WindowDuration:    time.Minute,                      // 1 minutes to reset the bucket
	BackendServerName: types.BackendServerName_EastAsia, // the current server
}
