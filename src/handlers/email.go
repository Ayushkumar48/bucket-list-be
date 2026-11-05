package handlers

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailSender struct {
	Host, Port, Username, Password string
}

func (s EmailSender) SendMagicLink(to, link string) error {
	if s.Username == "" || s.Password == "" {
		return fmt.Errorf("SMTP credentials are not configured")
	}
	if strings.ContainsAny(to+link, "\r\n") {
		return fmt.Errorf("invalid email or link")
	}
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	message := strings.Join([]string{
		"From: " + s.Username,
		"To: " + to,
		"Subject: Your Bucket List sign-in link",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		"Use this link to sign in to Bucket List:", link, "",
		"This link expires in 10 minutes and can be used only once.",
	}, "\r\n")
	return smtp.SendMail(s.Host+":"+s.Port, auth, s.Username, []string{to}, []byte(message))
}
