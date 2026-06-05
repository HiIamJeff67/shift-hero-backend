package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type SchedulePublicationStatus string

const (
	SchedulePublicationStatus_Draft     SchedulePublicationStatus = "Draft"
	SchedulePublicationStatus_Published SchedulePublicationStatus = "Published"
)

var AllSchedulePublicationStatuses = []SchedulePublicationStatus{
	SchedulePublicationStatus_Draft,
	SchedulePublicationStatus_Published,
}

var AllSchedulePublicationStatusStrings = []string{
	string(SchedulePublicationStatus_Draft),
	string(SchedulePublicationStatus_Published),
}

func (sps SchedulePublicationStatus) Name() string {
	return reflect.TypeOf(sps).Name()
}

func (sps *SchedulePublicationStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*sps = SchedulePublicationStatus(string(v))
		return nil
	case string:
		*sps = SchedulePublicationStatus(v)
		return nil
	}
	return scanError(value, sps)
}

func (sps SchedulePublicationStatus) Value() (driver.Value, error) {
	return string(sps), nil
}

func (sps SchedulePublicationStatus) String() string {
	return string(sps)
}

func (sps *SchedulePublicationStatus) IsValidEnum() bool {
	for _, enum := range AllSchedulePublicationStatuses {
		if *sps == enum {
			return true
		}
	}
	return false
}

func ConvertStringToSchedulePublicationStatus(enumString string) (*SchedulePublicationStatus, error) {
	for _, status := range AllSchedulePublicationStatuses {
		if string(status) == enumString {
			return &status, nil
		}
	}
	return nil, fmt.Errorf("invalid schedule publication status: %s", enumString)
}
