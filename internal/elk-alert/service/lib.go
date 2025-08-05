package service

import (
	"context"
	"elk-alert/internal/elk-alert/config"
	"elk-alert/internal/elk-alert/models"
	"elk-alert/internal/elk-alert/repositories"
	"log"
	"time"
)

type AlertHandlerService struct {
	cfg                  config.AlertHandlerConfig
	alertEventRepository repositories.AlertEventRepository
	senders              map[models.AlertChannel]repositories.AlertSender
}

func NewAlertHandlerService(
	cfg config.AlertHandlerConfig,
	alertEventRepository repositories.AlertEventRepository,
	senders ...repositories.AlertSender,
) *AlertHandlerService {
	senderMap := make(map[models.AlertChannel]repositories.AlertSender)
	for _, sender := range senders {
		senderMap[sender.GetAlertChannel()] = sender
	}

	return &AlertHandlerService{
		cfg:                  cfg,
		alertEventRepository: alertEventRepository,
		senders:              senderMap,
	}
}

func (s *AlertHandlerService) Start(ctx context.Context, index string) {
	for {
		log.Printf("Fetching alert event list %s...", index)
		alertEvents, err := s.alertEventRepository.GetAlertList(ctx, index)
		if err != nil {
			log.Printf("Failed to get alert event list %s: %v", index, err)
			time.Sleep(time.Second * 10)
			continue
		}
		log.Printf("alert event list fetched %s: %d", index, len(alertEvents))

		for _, alertEvent := range alertEvents {
			newTimestamp, err := time.Parse(time.RFC3339, alertEvent.Timestamp)
			if err != nil {
				log.Printf("Failed to parse timestamp %s: %v", index, err)
				continue
			}
			for _, alert := range alertEvent.Alerts {
				sender, exists := s.senders[alert.Channel]
				if !exists {
					log.Printf("Unknown alert channel %s: %s", index, alert.Channel)
					continue
				}

				newStatus := alertEvent.Status
				oldStatus, oldTimestamp, err := s.alertEventRepository.GetAlertStatus(ctx, alertEvent, alert.Channel)
				if err == nil {
					if oldTimestamp != nil && (oldTimestamp.After(newTimestamp) || oldTimestamp.Equal(newTimestamp)) {
						continue
					}
					if oldStatus == newStatus {
						continue
					}
				} else {
					log.Printf("Failed to get alert status %s: %v", index, err)
				}
				if err := s.alertEventRepository.SetAlertStatus(
					ctx, alertEvent, alert.Channel, sender.GetTTL(),
				); err != nil {
					log.Printf("Failed to set alert status %s: %v", index, err)
				}
				log.Printf("Sending %v to %s", alertEvent.AlertId, alert.Channel)
				_ = sender.Send(alert)
			}
		}

		time.Sleep(s.cfg.Interval)
	}
}
