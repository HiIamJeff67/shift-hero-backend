package emails

import (
	"fmt"
	"time"

	exceptions "github.com/your-org/go-start-monolithic-kit/app/exceptions"
	types "github.com/your-org/go-start-monolithic-kit/shared/types"
)

const ValidationEmailSubjectTemplate = "Verify Your Identity - %s Authentication Code"

var _validationEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/validation_email_template.html",
	DataMap:      map[string]any{},
}

func AsyncSendValidationEmail(
	to string,
	userName string,
	authCode string,
	userAgent string,
	expiredAt time.Time,
) *exceptions.Exception {
	remainingMinutes := int(time.Until(expiredAt).Minutes())

	dataMap := baseTemplateData()
	dataMap["UserName"] = userName
	dataMap["Email"] = to
	dataMap["AuthCode"] = authCode
	dataMap["UserAgent"] = userAgent
	dataMap["ExpiryMinutes"] = remainingMinutes
	dataMap["RequestTime"] = time.Now().Format("2006-01-02 15:04:05 MST")

	_validationEmailRenderer.DataMap = dataMap

	body, exception := _validationEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	emailObject := EmailObject{
		To:               to,
		Subject:          fmt.Sprintf(ValidationEmailSubjectTemplate, officialName),
		Body:             body,
		EmailContentType: types.EmailContentType_HTML,
	}

	exception = AppEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Validation, 3, 2)
	if exception != nil {
		return exception
	}

	return nil
}
