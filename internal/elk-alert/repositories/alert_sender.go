package repositories

import (
	"elk-alert/internal/elk-alert/models"
	"time"
)

type AlertSender interface {
	Send(alert models.Alert) error
	GetAlertChannel() models.AlertChannel
	GetTTL() time.Duration
}
