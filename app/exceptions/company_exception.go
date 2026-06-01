package exceptions

import (
	"fmt"
	"net/http"

	traces "github.com/HiIamJeff67/shift-hero-backend/app/monitor/traces"
)

const (
	_ExceptionBaseCode_Company ExceptionCode = CompanyExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	CompanyExceptionSubDomainCode ExceptionCode   = 52
	ExceptionBaseCode_Company     ExceptionCode   = _ExceptionBaseCode_Company + ReservedExceptionCode
	ExceptionPrefix_Company       ExceptionPrefix = "Company"
)

type CompanyExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	TypeExceptionDomain
}

var Company = &CompanyExceptionDomain{
	BaseCode: ExceptionBaseCode_Company,
	Prefix:   ExceptionPrefix_Company,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Company,
		_Prefix:   ExceptionPrefix_Company,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Company,
		_Prefix:   ExceptionPrefix_Company,
	},
}

func (d *CompanyExceptionDomain) Forbidden(message string) *Exception {
	if message == "" {
		message = "No permission to perform this company operation"
	}
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "Forbidden",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusForbidden,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *CompanyExceptionDomain) DuplicateMember(companyId string, userId string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "DuplicateMember",
		IsInternal:     false,
		Message:        fmt.Sprintf("User %s is already a member of company %s", userId, companyId),
		HTTPStatusCode: http.StatusConflict,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *CompanyExceptionDomain) BadRequest(message string) *Exception {
	if message == "" {
		message = "Invalid request"
	}
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "BadRequest",
		IsInternal:     false,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}
