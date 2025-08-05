package models

type AlertEvent struct {
	TemplateVersion string  `json:"template_version" mapstructure:"template_version"`
	AlertName       string  `json:"alert_name" mapstructure:"alert_name"`
	AlertId         string  `json:"alert_id" mapstructure:"alert_id"`
	Timestamp       string  `json:"timestamp" mapstructure:"timestamp"`
	AlertSource     string  `json:"alert_source" mapstructure:"alert_source"`
	Status          string  `json:"status" mapstructure:"status"`
	TeamName        string  `json:"team_name" mapstructure:"team_name"`
	Alerts          []Alert `json:"alerts" mapstructure:"alerts"`
}

type Alert struct {
	Channel     AlertChannel `json:"channel" mapstructure:"channel"`
	GroupName   string       `json:"group_name" mapstructure:"group_name"`
	Responsible []string     `json:"responsible" mapstructure:"responsible"`
	Timestamp   string       `json:"timestamp" mapstructure:"timestamp"`
	Message     AlertMessage `json:"message" mapstructure:"message"`
}

type AlertMessage struct {
	Title   string         `json:"title" mapstructure:"title"`
	Summary []Alertsummary `json:"summary" mapstructure:"summary"`
}

type Alertsummary struct {
	Label string `json:"label" mapstructure:"label"`
	Value string `json:"value" mapstructure:"value"`
}

type AlertChannel string

const SmsChannel AlertChannel = "sms"
const SlackChannel AlertChannel = "slack"
const TelegramChannel AlertChannel = "telegram"
const EmailChannel AlertChannel = "email"
