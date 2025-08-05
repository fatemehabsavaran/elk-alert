package repositories

import (
	"crypto/tls"
	"elk-alert/internal/elk-alert/models"
	gomail "gopkg.in/mail.v2"
	"log"
	"strings"
	"time"
)

type EmailSenderConfig struct {
	Host     string        `json:"host" mapstructure:"host"`
	Username string        `json:"username" mapstructure:"username"`
	Password string        `json:"password" mapstructure:"password"`
	Port     int           `json:"port" mapstructure:"port"`
	Sender   string        `json:"sender" mapstructure:"sender"`
	AlertTTL time.Duration `json:"alert_ttl" mapstructure:"alert_ttl"`
}

type EmailSender struct {
	cfg      EmailSenderConfig
	sendChan chan models.Alert
	dialer   *gomail.Dialer
}

func (k *EmailSender) GetTTL() time.Duration {
	return k.cfg.AlertTTL
}

func (k *EmailSender) GetAlertChannel() models.AlertChannel {
	return models.EmailChannel
}

func NewEmailSender(cfg EmailSenderConfig) *EmailSender {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Sender, cfg.Password)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return &EmailSender{
		cfg:      cfg,
		sendChan: make(chan models.Alert, 10),
		dialer:   d,
	}
}

func (k *EmailSender) Send(alert models.Alert) error {
	k.sendChan <- alert

	return nil
}

func (k *EmailSender) send(alert models.Alert) error {
	sb := strings.Builder{}
	sb.WriteString("Title: ")
	sb.WriteString(alert.Message.Title)
	sb.WriteString("\nsummary:")
	for _, summary := range alert.Message.Summary {
		sb.WriteString("\n")
		sb.WriteString(summary.Label)
		sb.WriteString(": ")
		sb.WriteString(summary.Value)
	}
	message := sb.String()

	m := gomail.NewMessage()
	m.SetHeader("From", k.cfg.Sender)
	m.SetHeader("To", alert.Responsible...)
	m.SetHeader("Subject", alert.Message.Title)
	m.SetBody("text/plain", message)

	log.Printf("Sending...")
	if err := k.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %s", err)
		return err
	}
	log.Printf("Done")
	return nil
}

func (k *EmailSender) Start() {
	go func() {
		for alert := range k.sendChan {
			err := k.send(alert)
			if err != nil {
				log.Printf("Failed to send sms: %s", err.Error())
			}
		}
	}()
}
