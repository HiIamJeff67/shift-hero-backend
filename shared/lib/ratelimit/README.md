# Rate Limit Library

## Overview

`shared/lib/ratelimit` provides rate-limit utilities used by middleware and interceptors:

- `WeakRateLimiter`: simple in-memory limiter.
- `HybridRateLimiter`: local token-bucket + Redis-synced global counters.
- `ReusableBufferPool`: `sync.Pool` wrapper for reusable `bytes.Buffer`.

## Key APIs

```go
type WeakRateLimiter
func NewWeakRateLimiter(requestsPerSecond int) *WeakRateLimiter
func (lb *WeakRateLimiter) Allow() bool

type HybridRateLimiter
func NewHybridRateLimiter(
	rateLimit rate.Limit,
	burst int,
	userLimit int32,
	windowDuration time.Duration,
	backendServerName types.BackendServerName,
	isAuthorizedLimiter bool,
) *HybridRateLimiter
func (hrl *HybridRateLimiter) AllowByFingerprint(fingerprint string) (bool, int32)
func (hrl *HybridRateLimiter) AllowByUserId(userId uuid.UUID) (bool, int32)
func (hrl *HybridRateLimiter) Allow(key string) (bool, int32)
func (hrl *HybridRateLimiter) AllowN(key string, now time.Time, n int) (bool, int32)
func (hrl *HybridRateLimiter) GetStatus() map[string]interface{}
func (hrl *HybridRateLimiter) GetDetailStatus() map[string]interface{}
func (hrl *HybridRateLimiter) Stop()

type ReusableBufferPool
func NewReusableBufferPool() *ReusableBufferPool
func (p *ReusableBufferPool) Get() *bytes.Buffer
func (p *ReusableBufferPool) Put(buffer *bytes.Buffer)
```

## Usage in This Project

- `app/middlewares/authorized_rate_limit_middleware.go`
- `app/middlewares/unauthorized_rate_limit_middleware.go`
- `app/middlewares/timeout_middleware.go`
- `app/interceptors/shareable_response_writer_interceptor.go`
