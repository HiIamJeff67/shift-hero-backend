package emails

import (
	"strings"

	util "github.com/HiIamJeff67/shift-hero-backend/app/util"
)

var (
	officialName = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_NAME", "Backend Starter"))
	officialMail = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_GMAIL", "noreply@example.com"))
	officialPass = util.GetEnv("APP_OFFICIAL_GOOGLE_APPLICATION_PASSWORD", "")
	smtpUsername = strings.TrimSpace(util.GetEnv("SMTP_USERNAME", officialMail))
	smtpPassword = util.GetEnv("SMTP_PASSWORD", officialPass)

	officialWebsiteURL         = strings.TrimSuffix(strings.TrimSpace(util.GetEnv("APP_OFFICIAL_WEBSITE_URL", "http://localhost:7777")), "/")
	officialHelpURL            = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_HELP_URL", officialWebsiteURL+"/help"))
	officialSupportURL         = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_SUPPORT_URL", officialWebsiteURL+"/support"))
	officialPrivacyURL         = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_PRIVACY_URL", officialWebsiteURL+"/privacy"))
	officialTermsURL           = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_TERMS_URL", officialWebsiteURL+"/terms"))
	officialSecurityURL        = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_SECURITY_URL", officialWebsiteURL+"/security"))
	officialSecurityReviewURL  = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_SECURITY_REVIEW_URL", officialWebsiteURL+"/security/review"))
	officialAccountSettingsURL = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_ACCOUNT_SETTINGS_URL", officialWebsiteURL+"/account/settings"))
	officialLoginURL           = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_LOGIN_URL", officialWebsiteURL+"/login"))

	officialSupportEmail  = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_SUPPORT_EMAIL", officialMail))
	officialSecurityEmail = strings.TrimSpace(util.GetEnv("APP_OFFICIAL_SECURITY_EMAIL", officialMail))
)

func baseTemplateData() map[string]any {
	return map[string]any{
		"OfficialName":               officialName,
		"OfficialWebsiteURL":         officialWebsiteURL,
		"OfficialHelpURL":            officialHelpURL,
		"OfficialSupportURL":         officialSupportURL,
		"OfficialPrivacyURL":         officialPrivacyURL,
		"OfficialTermsURL":           officialTermsURL,
		"OfficialSecurityURL":        officialSecurityURL,
		"OfficialSecurityReviewURL":  officialSecurityReviewURL,
		"OfficialAccountSettingsURL": officialAccountSettingsURL,
		"OfficialLoginURL":           officialLoginURL,
		"OfficialSupportMailTo":      "mailto:" + officialSupportEmail,
		"OfficialSecurityMailTo":     "mailto:" + officialSecurityEmail,
	}
}
