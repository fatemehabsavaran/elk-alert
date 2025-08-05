package repositories

import (
	"context"
	"elk-alert/internal/elk-alert/models"
	"time"
)

type AlertEventRepository interface {
	GetAlertList(ctx context.Context, index string) ([]models.AlertEvent, error)
	GetAlertStatus(ctx context.Context, alert models.AlertEvent, channel models.AlertChannel) (string, *time.Time, error)
	SetAlertStatus(ctx context.Context, alert models.AlertEvent, channel models.AlertChannel, duration time.Duration) error
}
