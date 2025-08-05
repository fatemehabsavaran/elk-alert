package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"net"
	"net/http"
	"time"
)

type ElasticConfig struct {
	ApiKey    string   `json:"api_key,omitempty" mapstructure:"api_key"`
	Addresses []string `json:"addresses,omitempty" mapstructure:"addresses"`
}

type Condition map[string]any

type ElasticConnector struct {
	client *elasticsearch.Client
	cfg    ElasticConfig
}

func NewElasticConnector(cfg ElasticConfig) (*ElasticConnector, error) {
	clientCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		APIKey:    cfg.ApiKey,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	client, err := elasticsearch.NewClient(clientCfg)
	if err != nil {
		return nil, err
	}

	return &ElasticConnector{
		cfg:    cfg,
		client: client,
	}, nil
}

func (e *ElasticConnector) GetClient() *elasticsearch.Client {
	return e.client
}
