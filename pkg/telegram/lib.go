package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type TelegramConfig struct {
	BotToken string        `json:"bot_token,omitempty" mapstructure:"bot_token"`
	Debug    bool          `json:"debug,omitempty" mapstructure:"debug"`
	AlertTTL time.Duration `json:"alert_ttl,omitempty" mapstructure:"alert_ttl"`
}

type TelegramConnector struct {
	cfg TelegramConfig
	bot *tgbotapi.BotAPI
}

func NewTelegramConnector(cfg TelegramConfig) (*TelegramConnector, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = cfg.Debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return &TelegramConnector{
		cfg: cfg,
		bot: bot,
	}, nil
}

func (t *TelegramConnector) GetBot() *tgbotapi.BotAPI {
	return t.bot
}

func (t *TelegramConnector) GetConfig() TelegramConfig {
	return t.cfg
}
