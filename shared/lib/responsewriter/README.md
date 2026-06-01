# ResponseWriter Library

## Overview

`shared/lib/responsewriter` wraps `gin.ResponseWriter` with a buffered writer.

It allows handlers/middlewares to:

- capture response body in memory,
- adjust headers/status after downstream handlers,
- flush to original writer at controlled timing.

## Key APIs

```go
type ResponseWriter
func NewResponseWriter(responseWriter gin.ResponseWriter, Body *bytes.Buffer) *ResponseWriter
func (rw *ResponseWriter) Header() http.Header
func (rw *ResponseWriter) WriteHeader(code int)
func (rw *ResponseWriter) Write(data []byte) (int, error)
func (rw *ResponseWriter) WriteString(s string) (int, error)
func (rw *ResponseWriter) WriteHeaderNow()
func (rw *ResponseWriter) Status() int
func (rw *ResponseWriter) Size() int
func (rw *ResponseWriter) FlushToOriginalWriter() error
func (rw *ResponseWriter) FreeBuffer()
```

## Usage in This Project

- `app/middlewares/timeout_middleware.go`
- `app/interceptors/embedded_interceptor.go`
- `app/interceptors/refresh_token_interceptor.go`
- `app/interceptors/shareable_response_writer_interceptor.go`
