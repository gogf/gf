package limiter

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	Tokens           = "tokens"         // Tokens represents the key for token count in the rate limiter bucket
	LastTime         = "last_time"      // LastTime represents the key for last update time in the rate limiter bucket
	DefaultRate      = 100              // DefaultRate is the default rate of token generation per second
	DefaultShards    = 16               // DefaultShards is the default number of shards for concurrent access
	DefaultCapacity  = 1024             // DefaultCapacity is the default capacity of the token bucket
	DefaultExpire    = 10 * time.Second // DefaultExpire is the default expiration time for cached entries
	DefaultKeyPrefix = "Rate-limiter:"
)

// DefaultKeyFunc is the default function to generate key from request, using client IP
func DefaultKeyFunc(r *ghttp.Request) string {
	return DefaultKeyPrefix + r.GetClientIp()
}

// DefaultAllowHandler is the default handler for allowed requests, continues to next middleware
func DefaultAllowHandler(r *ghttp.Request) {
	r.Middleware.Next()
}

// DefaultDenyHandler is the default handler for denied requests, returns 429 status
func DefaultDenyHandler(r *ghttp.Request) {
	r.Response.WriteStatus(http.StatusTooManyRequests, "Too Many Requests")
	r.ExitAll()
}
