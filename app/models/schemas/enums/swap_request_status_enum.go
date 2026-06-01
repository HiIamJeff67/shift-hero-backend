package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type SwapRequestStatus string

const (
	SwapRequestStatus_Open      SwapRequestStatus = "Open"
	SwapRequestStatus_Claimed   SwapRequestStatus = "Claimed"
	SwapRequestStatus_Approved  SwapRequestStatus = "Approved"
	SwapRequestStatus_Cancelled SwapRequestStatus = "Cancelled"
)

var AllSwapRequestStatuses = []SwapRequestStatus{
	SwapRequestStatus_Open,
	SwapRequestStatus_Claimed,
	SwapRequestStatus_Approved,
	SwapRequestStatus_Cancelled,
}

var AllSwapRequestStatusStrings = []string{
	string(SwapRequestStatus_Open),
	string(SwapRequestStatus_Claimed),
	string(SwapRequestStatus_Approved),
	string(SwapRequestStatus_Cancelled),
}

func (srs SwapRequestStatus) Name() string {
	return reflect.TypeOf(srs).Name()
}

func (srs *SwapRequestStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*srs = SwapRequestStatus(string(v))
		return nil
	case string:
		*srs = SwapRequestStatus(v)
		return nil
	}
	return scanError(value, srs)
}

func (srs SwapRequestStatus) Value() (driver.Value, error) {
	return string(srs), nil
}

func (srs SwapRequestStatus) String() string {
	return string(srs)
}

func (srs *SwapRequestStatus) IsValidEnum() bool {
	for _, enum := range AllSwapRequestStatuses {
		if *srs == enum {
			return true
		}
	}
	return false
}

func ConvertStringToSwapRequestStatus(enumString string) (*SwapRequestStatus, error) {
	for _, status := range AllSwapRequestStatuses {
		if string(status) == enumString {
			return &status, nil
		}
	}
	return nil, fmt.Errorf("invalid swap request status: %s", enumString)
}
