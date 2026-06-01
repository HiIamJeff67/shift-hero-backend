package emails

import (
	"fmt"
	"time"

	exceptions "github.com/your-org/go-start-monolithic-kit/app/exceptions"
	types "github.com/your-org/go-start-monolithic-kit/shared/types"
)

const SecurityAlertEmailSubjectTemplate = "%s Security Alert - Suspicious Activity Detected"

var _securityAlertEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/security_alert_email_template.html",
	DataMap:      map[string]any{},
}

func AsyncSendSecurityAlertEmail(
	to string,
	userName string,
	status string,
	alertType string,
	reason string,
	timeOfOccurrence time.Time,
	otherDetails string,
) *exceptions.Exception {
	dataMap := baseTemplateData()
	dataMap["UserName"] = userName
	dataMap["Status"] = status
	dataMap["AlertType"] = alertType
	dataMap["Reason"] = reason
	dataMap["TimeOfOccurrence"] = timeOfOccurrence
	dataMap["OtherDetails"] = otherDetails

	_securityAlertEmailRenderer.DataMap = dataMap

	body, exception := _securityAlertEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	emailObject := EmailObject{
		To:               to,
		Subject:          fmt.Sprintf(SecurityAlertEmailSubjectTemplate, officialName),
		Body:             body,
		EmailContentType: types.EmailContentType_HTML,
	}

	exception = AppEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Security, 3, 5)
	if exception != nil {
		return exception
	}

	return nil
}
