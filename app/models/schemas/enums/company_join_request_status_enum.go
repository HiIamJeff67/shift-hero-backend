package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type CompanyJoinRequestStatus string

const (
	CompanyJoinRequestStatus_Pending   CompanyJoinRequestStatus = "Pending"
	CompanyJoinRequestStatus_Approved  CompanyJoinRequestStatus = "Approved"
	CompanyJoinRequestStatus_Rejected  CompanyJoinRequestStatus = "Rejected"
	CompanyJoinRequestStatus_Cancelled CompanyJoinRequestStatus = "Cancelled"
)

var AllCompanyJoinRequestStatuses = []CompanyJoinRequestStatus{
	CompanyJoinRequestStatus_Pending,
	CompanyJoinRequestStatus_Approved,
	CompanyJoinRequestStatus_Rejected,
	CompanyJoinRequestStatus_Cancelled,
}

var AllCompanyJoinRequestStatusStrings = []string{
	string(CompanyJoinRequestStatus_Pending),
	string(CompanyJoinRequestStatus_Approved),
	string(CompanyJoinRequestStatus_Rejected),
	string(CompanyJoinRequestStatus_Cancelled),
}

func (cjrs CompanyJoinRequestStatus) Name() string {
	return reflect.TypeOf(cjrs).Name()
}

func (cjrs *CompanyJoinRequestStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*cjrs = CompanyJoinRequestStatus(string(v))
		return nil
	case string:
		*cjrs = CompanyJoinRequestStatus(v)
		return nil
	}
	return scanError(value, cjrs)
}

func (cjrs CompanyJoinRequestStatus) Value() (driver.Value, error) {
	return string(cjrs), nil
}

func (cjrs CompanyJoinRequestStatus) String() string {
	return string(cjrs)
}

func (cjrs *CompanyJoinRequestStatus) IsValidEnum() bool {
	for _, enum := range AllCompanyJoinRequestStatuses {
		if *cjrs == enum {
			return true
		}
	}
	return false
}

func ConvertStringToCompanyJoinRequestStatus(enumString string) (*CompanyJoinRequestStatus, error) {
	for _, status := range AllCompanyJoinRequestStatuses {
		if string(status) == enumString {
			return &status, nil
		}
	}
	return nil, fmt.Errorf("invalid company join request status: %s", enumString)
}
