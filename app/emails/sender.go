package emails

import (
	"errors"
	"strings"

	"gopkg.in/gomail.v2"

	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	util "github.com/HiIamJeff67/shift-hero-backend/app/util"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

/* ============================== Initialization & Instance ============================== */

type EmailSender struct {
	Host     string
	Port     int
	UserName string
	Password string
	From     string
}

var (
	AppEmailSender = &EmailSender{
		Host:     util.GetEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:     util.GetIntEnv("SMTP_PORT", 587),
		UserName: smtpUsername,
		Password: smtpPassword,
		From:     officialName + "<" + officialMail + ">",
	}
)

func (s *EmailSender) AsyncSend(to string, subject string, body string, contentType types.EmailContentType) *exceptions.Exception {
	if !contentType.IsValidEnum() {
		return exceptions.Email.InvalidEmailContentType(string(contentType))
	}
	if strings.TrimSpace(s.UserName) == "" || strings.EqualFold(strings.TrimSpace(s.UserName), "noreply@example.com") {
		return exceptions.Email.FailedToSendEmailWithSubject(subject).WithOrigin(
			errors.New("missing or placeholder SMTP username (set SMTP_USERNAME or APP_OFFICIAL_GMAIL)"),
		)
	}
	if strings.TrimSpace(s.Password) == "" {
		return exceptions.Email.FailedToSendEmailWithSubject(subject).WithOrigin(
			errors.New("missing SMTP password (set SMTP_PASSWORD or APP_OFFICIAL_GOOGLE_APPLICATION_PASSWORD)"),
		)
	}

	contentTypeString := contentType.String()

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody(contentTypeString, body)

	d := gomail.NewDialer(s.Host, s.Port, s.UserName, s.Password)
	if err := d.DialAndSend(m); err != nil {
		return exceptions.Email.FailedToSendEmailWithSubject(subject).WithOrigin(err)
	}
	return nil
}
