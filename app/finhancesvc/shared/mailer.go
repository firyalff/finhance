package shared

import (
	"net/smtp"
)

func SMTPAuth(username, password, hostURL string) smtp.Auth {
	return smtp.PlainAuth("", username, password, hostURL)
}

type SMTPMailConfig struct {
	Sender    string
	Recipient string
	Subject   string
	Body      string
}

func SMTPSendMail(smtpHostURL, smtpHostPort string, smtpSession smtp.Auth, smtpMail SMTPMailConfig) (err error) {
	content := "From: " + smtpMail.Sender + "\n" +
		"To: " + smtpMail.Recipient + "\n" +
		"Subject: " + smtpMail.Subject + "\n\n" +
		smtpMail.Body

	return smtp.SendMail(smtpHostURL+":"+smtpHostPort, smtpSession, smtpMail.Sender, []string{smtpMail.Recipient}, []byte(content))
}
