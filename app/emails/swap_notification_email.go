package emails

import (
	"fmt"

	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

func AsyncSendSwapClaimedEmail(to string, companyName string, assignmentSummary string) *exceptions.Exception {
	emailObject := EmailObject{
		To:               to,
		Subject:          fmt.Sprintf("[%s] Your swap request has been claimed", officialName),
		Body:             fmt.Sprintf("Your swap request in %s has been claimed.\nAssignment: %s", companyName, assignmentSummary),
		EmailContentType: types.EmailContentType_PlainText,
	}

	if exception := AppEmailWorkerManager.Enqueue(emailObject, EmailTaskType_SwapClaimed, 3, 4); exception != nil {
		return exception
	}
	return nil
}

func AsyncSendSwapApprovedEmail(to string, companyName string, assignmentSummary string) *exceptions.Exception {
	emailObject := EmailObject{
		To:               to,
		Subject:          fmt.Sprintf("[%s] Your swap request has been approved", officialName),
		Body:             fmt.Sprintf("Your swap request in %s has been approved.\nAssignment: %s", companyName, assignmentSummary),
		EmailContentType: types.EmailContentType_PlainText,
	}

	if exception := AppEmailWorkerManager.Enqueue(emailObject, EmailTaskType_SwapApproved, 3, 4); exception != nil {
		return exception
	}
	return nil
}
