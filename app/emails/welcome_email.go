package emails

import (
	"fmt"

	exceptions "github.com/your-org/go-start-monolithic-kit/app/exceptions"
	types "github.com/your-org/go-start-monolithic-kit/shared/types"
)

const WelcomeEmailSubjectTemplate = "Welcome to %s - Thanks for the Registration"

var _welcomeEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/welcome_email_template.html",
	DataMap:      map[string]any{},
}

func AsyncSendWelcomeEmail(
	to string,
	userName string,
	status string,
) *exceptions.Exception {
	dataMap := baseTemplateData()
	dataMap["UserName"] = userName
	dataMap["Email"] = to
	dataMap["Status"] = status

	_welcomeEmailRenderer.DataMap = dataMap
	body, exception := _welcomeEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	emailObject := EmailObject{
		To:               to,
		Subject:          fmt.Sprintf(WelcomeEmailSubjectTemplate, officialName),
		Body:             body,
		EmailContentType: types.EmailContentType_HTML,
	}

	exception = AppEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Welcome, 3, 1)
	if exception != nil {
		return exception
	}

	return nil
}
