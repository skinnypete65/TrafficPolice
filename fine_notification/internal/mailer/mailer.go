package mailer

import (
	"fine_notification/internal/transport/dto"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

const (
	contentTypeText = "text/plain"
)

type Mailer struct {
	mailDialer *gomail.Dialer
}

func NewMailer(mailDialer *gomail.Dialer) *Mailer {
	return &Mailer{mailDialer: mailDialer}
}

func (m *Mailer) SendFineMessage(email Email, caseInfo dto.Case) error {

	gm := gomail.NewMessage()
	gm.SetHeader("From", email.From)
	gm.SetHeader("To", email.To)
	gm.SetHeader("Subject", email.Subject)
	gm.SetBody(contentTypeText, m.getBodyFromCase(caseInfo))
	log.Println(email, caseInfo)

	return m.mailDialer.DialAndSend(gm)
}

func (m *Mailer) getBodyFromCase(c dto.Case) string {
	return fmt.Sprintf(
		`Координаты места происшествия: %f,%f
Правонарушение: %s, значение: %s
Размер штрафа: %d
Дата: %s`,
		c.Camera.Latitude, c.Camera.Longitude,
		c.Violation.Name, c.ViolationValue,
		c.Violation.FineAmount,
		c.Date.Format(time.RFC850))
}
