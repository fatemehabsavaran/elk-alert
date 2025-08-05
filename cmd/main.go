package main

import (
	"context"
	appcfg "elk-alert/config"
	elkConfig "elk-alert/internal/elk-alert/config"
	"elk-alert/internal/elk-alert/service"
	"elk-alert/internal/repositories"
	"elk-alert/pkg/config"
	"elk-alert/pkg/elastic"
	"elk-alert/pkg/redis"
	"elk-alert/pkg/telegram"
	"fmt"
	"log"
	"net/http"
)

var Config appcfg.Config

func init() {
	Config = *config.Load[appcfg.Config]()
	Config.Elastic.ApiKey = config.GetEnv("ElasticApiKey", Config.Elastic.ApiKey)
	Config.Telegram.BotToken = config.GetEnv("TelegramBotToken", Config.Telegram.BotToken)
	Config.Sms.ApiKey = config.GetEnv("KavenegarApiKey", Config.Sms.ApiKey)
	Config.Email.Password = config.GetEnv("SMTPPassword", Config.Email.Password)
	log.Printf("%+v", Config)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	elasticClient, err := elastic.NewElasticConnector(Config.Elastic)
	if err != nil {
		log.Fatalf("Failed to create elastic connector: %v", err)
		return
	}

	tgClient, err := telegram.NewTelegramConnector(Config.Telegram)
	if err != nil {
		log.Fatalf("Failed to create telegram connector: %v", err)
		return
	}

	smsSender := repositories.NewSmsSender(Config.Sms)
	smsSender.Start()

	emailSender := repositories.NewEmailSender(Config.Email)
	emailSender.Start()

	tgSender := repositories.NewTelegramSender(tgClient)
	tgSender.Start()

	for _, alert := range Config.AlertHandler.Alerts {
		go func(alert elkConfig.AlertHandlerAlert) {
			redisClient := redis.NewRedisConnector(redis.RedisConfig{
				Addr:       Config.RedisConfig.Addr,
				Password:   Config.RedisConfig.Password,
				DB:         alert.RedisDB,
				PoolSize:   Config.RedisConfig.PoolSize,
				Timeout:    Config.RedisConfig.Timeout,
				ClientName: Config.RedisConfig.ClientName,
			})

			slackSender := repositories.NewSlackSender(repositories.SlackSenderConfig{
				WebhookURL: alert.SlackWebhook,
				AlertTTL:   Config.Slack.AlertTTL,
			})
			slackSender.Start()

			alertHandler := service.NewAlertHandlerService(
				Config.AlertHandler,
				repositories.NewAlertEventProvider(elasticClient, redisClient),
				smsSender,
				tgSender,
				slackSender,
				emailSender,
			)

			alertHandler.Start(context.Background(), alert.ElasticIndex)
		}(alert)
	}

	http.HandleFunc("/health", healthCheckHandler)

	port := ":8080"
	fmt.Printf("Starting healthcheck on %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
