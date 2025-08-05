package config

import (
	"elk-alert/internal/elk-alert/config"
	"elk-alert/internal/repositories"
	"elk-alert/pkg/elastic"
	"elk-alert/pkg/redis"
	"elk-alert/pkg/telegram"
)

type Config struct {
	Telegram     telegram.TelegramConfig        `json:"telegram,omitempty" mapstructure:"telegram"`
	Slack        repositories.SlackSenderConfig `json:"slack,omitempty" mapstructure:"slack"`
	Elastic      elastic.ElasticConfig          `json:"elastic,omitempty" mapstructure:"elastic"`
	Sms          repositories.SmsSenderConfig   `json:"sms,omitempty" mapstructure:"sms"`
	Email        repositories.EmailSenderConfig `json:"email,omitempty" mapstructure:"email"`
	AlertHandler config.AlertHandlerConfig      `json:"alert_handler,omitempty" mapstructure:"alert_handler"`
	RedisConfig  redis.RedisConfig              `json:"redis,omitempty" mapstructure:"redis"`
}
