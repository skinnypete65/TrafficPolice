package mailer

import (
	"fine_notification/internal/transport/dto"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"time"
)

const (
	contextTypeHtml = "text/html"
	fieldFrom       = "From"
	fieldTo         = "To"
	fieldSubject    = "subject"
	violationPrefix = "violation"
)

type Mailer struct {
	mailDialer *gomail.Dialer
}

func NewMailer(mailDialer *gomail.Dialer) *Mailer {
	return &Mailer{mailDialer: mailDialer}
}

func (m *Mailer) SendFineNotification(email Email, caseInfo dto.CaseWithImage) error {
	body := m.getBodyFromCase(caseInfo)

	gm := gomail.NewMessage()
	gm.SetHeader(fieldFrom, email.From)
	gm.SetHeader(fieldTo, email.To)
	gm.SetHeader(fieldSubject, email.Subject)
	gm.SetBody(contextTypeHtml, body)
	gm.Attach(
		fmt.Sprintf("%s.%s", violationPrefix, caseInfo.ImageExtension),
		gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(caseInfo.Image)
			return err
		}),
	)

	return m.mailDialer.DialAndSend(gm)
}

func (m *Mailer) getBodyFromCase(caseInfo dto.CaseWithImage) string {
	coordinatesInfo := fmt.Sprintf("Координаты места происшествия: %f,%f",
		caseInfo.Case.Camera.Latitude, caseInfo.Case.Camera.Longitude)
	violationInfo := fmt.Sprintf("Правонарушение: %s, значение: %s",
		caseInfo.Case.Violation.Name, caseInfo.Case.ViolationValue)
	fineAmountInfo := fmt.Sprintf("Размер штрафа: %d", caseInfo.Case.Violation.FineAmount)
	dateInfo := fmt.Sprintf("Дата: %s", caseInfo.Case.Date.Format(time.RFC850))

	return fmt.Sprintf(`
<p>%s</p>
<p>%s</p>
<p>%s</p>
<p>%s</p>`,
		coordinatesInfo, violationInfo, fineAmountInfo, dateInfo)
}
