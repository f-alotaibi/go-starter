package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/wneessen/go-mail"
)

func NewMailClient() (*mail.Client, error) {
	port, err := strconv.Atoi(os.Getenv("MAIL_SMTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("mail: port is not a valid integer")
	}

	client, err := mail.NewClient(os.Getenv("MAIL_SMTP_SERVER"), mail.WithPort(port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(os.Getenv("MAIL_SMTP_USER")), mail.WithPassword(os.Getenv("MAIL_SMTP_PASS")))
	if err != nil {
		return nil, err
	}

	return client, err
}
