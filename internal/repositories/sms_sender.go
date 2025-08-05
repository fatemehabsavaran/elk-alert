package repositories

import (
	"elk-alert/internal/elk-alert/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SmsSenderConfig struct {
	ApiKey   string        `json:"api_key,omitempty" mapstructure:"api_key"`
	Sender   string        `json:"sender,omitempty" mapstructure:"sender"`
	AlertTTL time.Duration `json:"alert_ttl,omitempty" mapstructure:"alert_ttl"`
}

type SmsSender struct {
	cfg        SmsSenderConfig
	sendChan   chan models.Alert
	httpClient *http.Client
}

func (k *SmsSender) GetTTL() time.Duration {
	return k.cfg.AlertTTL
}

func (k *SmsSender) GetAlertChannel() models.AlertChannel {
	return models.SmsChannel
}

func NewSmsSender(cfg SmsSenderConfig) *SmsSender {
	return &SmsSender{
		cfg:        cfg,
		sendChan:   make(chan models.Alert, 10),
		httpClient: http.DefaultClient,
	}
}

func (k *SmsSender) Send(alert models.Alert) error {
	k.sendChan <- alert

	return nil
}

func (k *SmsSender) send(alert models.Alert) error {
	apiURL := fmt.Sprintf("https://api.kavenegar.com/v1/%s/sms/send.json", k.cfg.ApiKey)

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

	receptorSB := strings.Builder{}
	for i, receiver := range alert.Responsible {
		receptorSB.WriteString(receiver)
		if i+1 != len(alert.Responsible) {
			receptorSB.WriteString(",")
		}
	}
	receiver := receptorSB.String()

	data := url.Values{}
	data.Set("receptor", receiver)
	data.Set("message", message)
	data.Set("sender", k.cfg.Sender)
	log.Printf("Sending from %s to %s: %s", k.cfg.Sender, receiver, message)
	log.Printf("To Url %s", apiURL)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		log.Printf("Failed to send request: %d %s", resp.StatusCode, string(body))
		return fmt.Errorf("failed to send request: %d %s", resp.StatusCode, string(body))
	}

	var result kavenegarResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return err
	}

	return nil
}

func (k *SmsSender) Start() {
	go func() {
		for alert := range k.sendChan {
			err := k.send(alert)
			if err != nil {
				log.Printf("Failed to send sms: %s", err.Error())
			}
		}
	}()
}

type kavenegarResponse struct {
	Entries []struct {
		MessageID int    `json:"messageid"`
		Sender    string `json:"sender"`
		Message   string `json:"message"`
		Date      int64  `json:"date"`
		Receptor  string `json:"receptor"`
	} `json:"entries"`
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
}
