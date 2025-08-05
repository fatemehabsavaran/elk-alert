package config

import "time"

type AlertHandlerConfig struct {
	Interval     time.Duration       `json:"interval,omitempty" mapstructure:"interval"`
	MarkDuration time.Duration       `json:"markDuration,omitempty" mapstructure:"markDuration"`
	Alerts       []AlertHandlerAlert `json:"alerts" mapstructure:"alerts"`
}

type AlertHandlerAlert struct {
	ElasticIndex string `json:"elasticIndex,omitempty" mapstructure:"elasticIndex"`
	RedisDB      int    `json:"redisDB,omitempty" mapstructure:"redisDB"`
	SlackWebhook string `json:"slackWebhook,omitempty" mapstructure:"slackWebhook"`
}
