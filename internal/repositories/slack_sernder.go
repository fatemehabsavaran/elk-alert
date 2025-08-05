package repositories

import (
	"bytes"
	"elk-alert/internal/elk-alert/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SlackSenderConfig struct {
	WebhookURL string        `json:"webhook_url,omitempty"`
	AlertTTL   time.Duration `json:"alert_ttl,omitempty" mapstructure:"alert_ttl"`
}

type SlackSender struct {
	cfg        SlackSenderConfig
	sendChan   chan models.Alert
	httpClient *http.Client
}

func NewSlackSender(cfg SlackSenderConfig) *SlackSender {
	return &SlackSender{
		cfg:        cfg,
		sendChan:   make(chan models.Alert, 10),
		httpClient: http.DefaultClient,
	}
}

func (s *SlackSender) GetTTL() time.Duration {
	return s.cfg.AlertTTL
}

func (s *SlackSender) GetAlertChannel() models.AlertChannel {
	return models.SlackChannel
}

func (s *SlackSender) Send(alert models.Alert) error {
	s.sendChan <- alert
	return nil
}

func (s *SlackSender) send(alert models.Alert) error {
	message := struct {
		Text string `json:"text"`
	}{
		Text: s.formatMessage(alert),
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", s.cfg.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send request, status code: %d", resp.StatusCode)
		return err
	}

	return nil
}

func (s *SlackSender) formatMessage(alert models.Alert) string {
	var sb bytes.Buffer
	sb.WriteString("*Title:* ")
	sb.WriteString(alert.Message.Title)
	sb.WriteString("\n*Summary:*\n")
	for _, summary := range alert.Message.Summary {
		sb.WriteString("- ")
		sb.WriteString(summary.Label)
		sb.WriteString(": ")
		sb.WriteString(summary.Value)
		sb.WriteString("\n")
	}
	sb.WriteString("*Responsible:* ")
	for _, person := range alert.Responsible {
		sb.WriteString(person + " ")
	}
	return sb.String()
}

func (s *SlackSender) Start() {
	go func() {
		for alert := range s.sendChan {
			if err := s.send(alert); err != nil {
				log.Printf("Failed to send Slack message: %s", err.Error())
			}
		}
	}()
}
