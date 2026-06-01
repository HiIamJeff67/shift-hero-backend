package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type EmployeeRole string

const (
	EmployeeRole_Manager EmployeeRole = "Manager"
	EmployeeRole_Staff   EmployeeRole = "Staff"
)

var AllEmployeeRoles = []EmployeeRole{
	EmployeeRole_Manager,
	EmployeeRole_Staff,
}
var AllEmployeeRoleStrings = []string{
	string(EmployeeRole_Manager),
	string(EmployeeRole_Staff),
}

func (er EmployeeRole) Name() string {
	return reflect.TypeOf(er).Name()
}

func (er *EmployeeRole) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*er = EmployeeRole(string(v))
		return nil
	case string:
		*er = EmployeeRole(v)
		return nil
	}
	return scanError(value, er)
}

func (er EmployeeRole) Value() (driver.Value, error) {
	return string(er), nil
}

func (er EmployeeRole) String() string {
	return string(er)
}

func (er *EmployeeRole) IsValidEnum() bool {
	for _, enum := range AllEmployeeRoles {
		if *er == enum {
			return true
		}
	}
	return false
}

func ConvertStringToEmployeeRole(enumString string) (*EmployeeRole, error) {
	for _, employeeRole := range AllEmployeeRoles {
		if string(employeeRole) == enumString {
			return &employeeRole, nil
		}
	}
	return nil, fmt.Errorf("invalid employee roles: %s", enumString)
}
