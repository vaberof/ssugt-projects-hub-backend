package email

import (
	"bytes"
	"encoding/base64"

	"fmt"
	"log/slog"
	"net/smtp"
	"ssugt-projects-hub/pkg/logging/logs"
)

type Service interface {
	SendConfirmationEmail(to, code string) error
	SendEmail(to []string, subject, body string) error
}

type SmtpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewSmtpConfig() SmtpConfig {
	return SmtpConfig{
		Host:     "smtp.yandex.ru",
		Port:     587,
		Username: "razrabotkasgugit",
		Password: "cfwbugarreotampl",
		From:     "razrabotkasgugit@yandex.ru",
	}
}

type serviceImpl struct {
	log    *slog.Logger
	config SmtpConfig
}

func NewService(log *logs.Logs, config SmtpConfig) Service {
	return &serviceImpl{
		log:    log.WithName("email-service"),
		config: config,
	}
}

func (s *serviceImpl) SendConfirmationEmail(to, code string) error {
	subject, body := getVerificationMail(code)

	s.log.Debug("sending email", "to", to, "subject", subject, "body", body)

	if err := s.SendEmail([]string{to}, subject, body); err != nil {
		return fmt.Errorf("ошибка отправки email: %w", err)
	}

	s.log.Debug("successfully sent email", "to", to, "subject", subject, "body", body)

	return nil
}

func (s *serviceImpl) SendEmail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	subjectEncoded := fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))

	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", s.config.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to[0]))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subjectEncoded))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return smtp.SendMail(addr, auth, s.config.From, to, msg.Bytes())
}

func getVerificationMail(code string) (subject string, body string) {
	subject = "Подтверждение регистрации"
	body = fmt.Sprintf("Здравствуйте!\n\nВаш код подтверждения: %s\n\nЕсли вы не запрашивали код, проигнорируйте это письмо.", code)
	return
}
