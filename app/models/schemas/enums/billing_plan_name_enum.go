package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type BillingPlanName string

// IMPORTANT FOR DEVELOPERS:
// The billing plan names below are template placeholders only.
// In real deployments, these values MUST match your actual billing provider
// plan/product naming exactly (for example, names configured in PayPal/Stripe).
// If you change external billing plan names, update this enum + seed examples
// together to keep upgrade/payment logic consistent.
const (
	BillingPlanName_MonthlyFreePlan       BillingPlanName = "Template Monthly Free Plan"
	BillingPlanName_MonthlyProPlan        BillingPlanName = "Template Monthly Pro Plan"
	BillingPlanName_YearlyProPlan         BillingPlanName = "Template Yearly Pro Plan"
	BillingPlanName_MonthlyPremiumPlan    BillingPlanName = "Template Monthly Premium Plan"
	BillingPlanName_YearlyPremiumPlan     BillingPlanName = "Template Yearly Premium Plan"
	BillingPlanName_MonthlyUltimatePlan   BillingPlanName = "Template Monthly Ultimate Plan"
	BillingPlanName_YearlyUltimatePlan    BillingPlanName = "Template Yearly Ultimate Plan"
	BillingPlanName_MonthlyEnterprisePlan BillingPlanName = "Template Monthly Enterprise Plan"
	BillingPlanName_YearlyEnterprisePlan  BillingPlanName = "Template Yearly Enterprise Plan"
)

var AllBillingPlanNames = []BillingPlanName{
	BillingPlanName_MonthlyFreePlan,
	BillingPlanName_MonthlyProPlan,
	BillingPlanName_YearlyProPlan,
	BillingPlanName_MonthlyPremiumPlan,
	BillingPlanName_YearlyPremiumPlan,
	BillingPlanName_MonthlyUltimatePlan,
	BillingPlanName_YearlyUltimatePlan,
	BillingPlanName_MonthlyEnterprisePlan,
	BillingPlanName_YearlyEnterprisePlan,
}

var AllBillingPlanNameStrings = []string{
	string(BillingPlanName_MonthlyFreePlan),
	string(BillingPlanName_MonthlyProPlan),
	string(BillingPlanName_YearlyProPlan),
	string(BillingPlanName_MonthlyPremiumPlan),
	string(BillingPlanName_YearlyPremiumPlan),
	string(BillingPlanName_MonthlyUltimatePlan),
	string(BillingPlanName_YearlyUltimatePlan),
	string(BillingPlanName_MonthlyEnterprisePlan),
	string(BillingPlanName_YearlyEnterprisePlan),
}

func (bpn BillingPlanName) Name() string {
	return reflect.TypeOf(bpn).Name()
}

func (bpn *BillingPlanName) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*bpn = BillingPlanName(string(v))
		return nil
	case string:
		*bpn = BillingPlanName(v)
		return nil
	}
	return scanError(value, bpn)
}

func (bpn BillingPlanName) Value() (driver.Value, error) {
	return string(bpn), nil
}

func (bpn BillingPlanName) String() string {
	return string(bpn)
}

func (bpn *BillingPlanName) IsValidEnum() bool {
	for _, enum := range AllBillingPlanNames {
		if *bpn == enum {
			return true
		}
	}
	return false
}

func ConvertStringToBillingPlanName(enumString string) (*BillingPlanName, error) {
	for _, billingPlanName := range AllBillingPlanNames {
		if string(billingPlanName) == enumString {
			return &billingPlanName, nil
		}
	}
	return nil, fmt.Errorf("invalid billing plan name: %s", enumString)
}
