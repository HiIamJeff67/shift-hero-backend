package types

type TableName string

const (
	// public tables(accessable and mutatable by the client user and admin)
	TableName_UserTable        TableName = "UserTable"
	TableName_UserAccountTable TableName = "UserAccountTable"
	TableName_UserInfoTable    TableName = "UserInfoTable"
	TableName_UserSettingTable TableName = "UserSettingTable"

	TableName_UsersToBillingPlansTable TableName = "UsersToBillingPlansTable"

	// private tables(accessable by the client user and admin, but only mutatable by the admin)
	TableName_PlanLimitationTable TableName = "PlanLimitationTable"
	TableName_BillingPlanTable    TableName = "BillingPlanTable"
)

var _validTableNames = map[string]TableName{
	// public tables
	"UserTable":        TableName_UserTable,
	"UserAccountTable": TableName_UserAccountTable,
	"UserInfoTable":    TableName_UserInfoTable,
	"UserSettingTable": TableName_UserSettingTable,

	"UsersToBillingPlansTable": TableName_UsersToBillingPlansTable,

	// private tables
	"PlanLimitationTable": TableName_PlanLimitationTable,
	"BillingPlanTable":    TableName_BillingPlanTable,
}

func (tn TableName) String() string {
	return string(tn)
}

func IsTableName(tableName string) bool {
	_, ok := _validTableNames[tableName]
	return ok
}
func ConvertToTableName(tableName string) (TableName, bool) {
	validTableName, ok := _validTableNames[tableName]
	return validTableName, ok
}
