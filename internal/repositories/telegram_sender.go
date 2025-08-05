package repositories

import (
	"elk-alert/internal/elk-alert/models"
	"elk-alert/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

type TelegramSender struct {
	connector *telegram.TelegramConnector
	sendChan  chan models.Alert
}

func NewTelegramSender(connector *telegram.TelegramConnector) *TelegramSender {
	return &TelegramSender{
		connector: connector,
		sendChan:  make(chan models.Alert, 10),
	}
}

func (t *TelegramSender) GetTTL() time.Duration {
	return t.connector.GetConfig().AlertTTL
}

func (t *TelegramSender) GetAlertChannel() models.AlertChannel {
	return models.TelegramChannel
}

func (t *TelegramSender) Send(alert models.Alert) error {
	t.sendChan <- alert
	return nil
}

func (t *TelegramSender) send(alert models.Alert) error {
	chatId, err := strconv.Atoi(alert.GroupName)
	if err != nil {
		log.Printf("Error parsing chat ID from GroupName '%s': %s", alert.GroupName, err.Error())
		return err
	}

	sb := new(strings.Builder)
	sb.WriteString("Title: ")
	sb.WriteString(alert.Message.Title)
	sb.WriteString("\nSummary:\n")

	for _, summary := range alert.Message.Summary {
		sb.WriteString(summary.Label)
		sb.WriteString(": ")
		sb.WriteString(summary.Value)
		sb.WriteString("\n")
	}

	sb.WriteString("Responsible: ")
	sb.WriteString("\n")
	for _, person := range alert.Responsible {
		sb.WriteString(person)
		sb.WriteString(" ")
	}
	sb.WriteString("\n")

	tgMsg := tgbotapi.NewMessage(int64(chatId), sb.String())

	_, err = t.connector.GetBot().Send(tgMsg)
	if err != nil {
		log.Printf("Failed to send Telegram message: %s", err.Error())
		return err
	}

	return nil
}

func (t *TelegramSender) Start() {
	go func() {
		for alert := range t.sendChan {
			err := t.send(alert)
			if err != nil {
				log.Printf("Failed to send alert: %s", err.Error())
			}
		}
	}()
}
