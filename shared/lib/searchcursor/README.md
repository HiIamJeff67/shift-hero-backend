# Search Cursor Library

## Overview

`shared/lib/searchcursor` provides generic cursor encode/decode helpers for cursor-based pagination.

It serializes cursor fields to JSON and then Base64-encodes them.

## Key APIs

```go
type SearchCursor[SearchCursorFieldType any] struct {
	Fields SearchCursorFieldType
}

func New[SearchCursorFieldType any](fields SearchCursorFieldType) *SearchCursor[SearchCursorFieldType]
func (sc *SearchCursor[SearchCursorFieldType]) Encode() (*string, error)
func Decode[SearchCursorFieldType any](encoded string) (*SearchCursor[SearchCursorFieldType], error)
func EncodeFromData[SearchCursorType any](data SearchCursorType) (*string, error)
```

## Usage in This Project

- `app/services/user_service.go`
