package exceptions

import (
	"net/http"
	"time"

	traces "github.com/HiIamJeff67/shift-hero-backend/app/monitor/traces"
)

const (
	_ExceptionBaseCode_Scheduling ExceptionCode = SchedulingExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	SchedulingExceptionSubDomainCode ExceptionCode   = 53
	ExceptionBaseCode_Scheduling     ExceptionCode   = _ExceptionBaseCode_Scheduling + ReservedExceptionCode
	ExceptionPrefix_Scheduling       ExceptionPrefix = "Scheduling"
)

type SchedulingExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	TypeExceptionDomain
}

var Scheduling = &SchedulingExceptionDomain{
	BaseCode: ExceptionBaseCode_Scheduling,
	Prefix:   ExceptionPrefix_Scheduling,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Scheduling,
		_Prefix:   ExceptionPrefix_Scheduling,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Scheduling,
		_Prefix:   ExceptionPrefix_Scheduling,
	},
}

func (d *SchedulingExceptionDomain) InvalidTimeRange() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "InvalidTimeRange",
		IsInternal:     false,
		Message:        "EndAt must be later than StartAt",
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *SchedulingExceptionDomain) Forbidden(message string) *Exception {
	if message == "" {
		message = "No permission to perform this scheduling operation"
	}
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "Forbidden",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusForbidden,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *SchedulingExceptionDomain) InvalidSwapState(message string) *Exception {
	if message == "" {
		message = "Invalid swap request state transition"
	}
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "InvalidSwapState",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *SchedulingExceptionDomain) BadRequest(message string) *Exception {
	if message == "" {
		message = "Invalid request"
	}
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "BadRequest",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *SchedulingExceptionDomain) AIUnavailable() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         "AIUnavailable",
		IsInternal:     false,
		Message:        "The AI schedule analyst is temporarily unavailable",
		HTTPStatusCode: http.StatusServiceUnavailable,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *SchedulingExceptionDomain) AIGenerationFailed() *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         "AIGenerationFailed",
		IsInternal:     false,
		Message:        "Failed to generate schedule insights",
		HTTPStatusCode: http.StatusBadGateway,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *SchedulingExceptionDomain) InsightRangeTooLarge() *Exception {
	return &Exception{
		Code:           d.BaseCode + 7,
		Prefix:         d.Prefix,
		Reason:         "InsightRangeTooLarge",
		IsInternal:     false,
		Message:        "Schedule insight date range cannot exceed 31 days",
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *SchedulingExceptionDomain) AIUsageLimitExceeded(used int32, limit int32, resetAt time.Time) *Exception {
	return &Exception{
		Code:           d.BaseCode + 8,
		Prefix:         d.Prefix,
		Reason:         "AIUsageLimitExceeded",
		IsInternal:     false,
		Message:        "The monthly AI generation limit has been reached",
		HTTPStatusCode: http.StatusTooManyRequests,
		Details: map[string]any{
			"used":      used,
			"limit":     limit,
			"remaining": 0,
			"resetAt":   resetAt,
		},
		LastTrace: traces.GetTrace(1),
	}
}
