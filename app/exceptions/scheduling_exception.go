package exceptions

import (
	"net/http"

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
