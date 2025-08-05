package repositories

import (
  "context"
  "elk-alert/internal/elk-alert/models"
  "elk-alert/pkg/elastic"
  "elk-alert/pkg/redis"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/mitchellh/mapstructure"
  redis2 "github.com/redis/go-redis/v9"
  "io"
  "log"
  "strings"
  "time"
)

type AlertEventProvider struct {
  esConnector *elastic.ElasticConnector
  rConnector  *redis.RedisConnector
}

func NewAlertEventProvider(
  esConnector *elastic.ElasticConnector,
  rConnector *redis.RedisConnector,
) *AlertEventProvider {
  return &AlertEventProvider{
    esConnector: esConnector,
    rConnector:  rConnector,
  }
}

func (e *AlertEventProvider) GetAlertList(ctx context.Context, index string) ([]models.AlertEvent, error) {
  query := elastic.Condition{
    "query": elastic.Condition{
      "bool": elastic.Condition{
        "must": []elastic.Condition{
          {
            "range": elastic.Condition{
              "timestamp": elastic.Condition{
                "gte": "now-5m",
                "lte": "now",
              },
            },
          },
        },
      },
    },
  }

  queryBody, err := json.Marshal(query)
  if err != nil {
    return nil, err
  }

  res, err := e.esConnector.GetClient().Search(
    e.esConnector.GetClient().Search.WithContext(ctx),
    e.esConnector.GetClient().Search.WithIndex(index),
    e.esConnector.GetClient().Search.WithBody(strings.NewReader(string(queryBody))),
    e.esConnector.GetClient().Search.WithTrackTotalHits(true),
    e.esConnector.GetClient().Search.WithRequestCache(false),
    e.esConnector.GetClient().Search.WithPretty(),
  )

  if err != nil {
    return nil, err
  }
  defer res.Body.Close()

  if res.IsError() {
    return nil, fmt.Errorf(res.Status())
  }

  var esResp struct {
    Hits struct {
      Hits []struct {
        SourceMap map[string]any `json:"_source"`
      } `json:"hits"`
    } `json:"hits"`
  }

  body, err := io.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  if err := json.Unmarshal(body, &esResp); err != nil {
    return nil, err
  }

  alertEvents := make([]models.AlertEvent, 0)
  for _, hit := range esResp.Hits.Hits {
    var alertEvent models.AlertEvent
    if err = mapstructure.Decode(hit.SourceMap, &alertEvent); err != nil {
      log.Printf("Unknown alert event recode: %+v", hit.SourceMap)
      continue
    }
    log.Printf("Alert event: %+v", alertEvent)

    alertEvents = append(alertEvents, alertEvent)
  }

  return alertEvents, nil
}

func getAlertKey(alertId string, channel models.AlertChannel) string {
  return fmt.Sprintf("alert_%s_%s", alertId, channel)
}

func getAlertValue(alert models.AlertEvent) string {
  return fmt.Sprintf("%s_%s", alert.Status, alert.Timestamp)
}

func extractAlertValue(value string) (string, time.Time) {
  splitStr := strings.SplitN(value, "_", 2)

  t, err := time.Parse(time.RFC3339, splitStr[1])
  if err != nil {
    t = time.UnixMilli(0)
  }

  return splitStr[0], t
}

func (e *AlertEventProvider) GetAlertStatus(
  ctx context.Context,
  alert models.AlertEvent,
  channel models.AlertChannel,
) (string, *time.Time, error) {
  res, err := e.rConnector.GetClient().Get(ctx, getAlertKey(alert.AlertId, channel)).Result()
  if err != nil {
    if errors.Is(err, redis2.Nil) {
      return "", nil, nil
    }

    return "", nil, err
  }

  oldStatus, oldTimestamp := extractAlertValue(res)

  return oldStatus, &oldTimestamp, nil
}

func (e *AlertEventProvider) SetAlertStatus(
  ctx context.Context,
  alert models.AlertEvent,
  channel models.AlertChannel,
  duration time.Duration,
) error {
  _, err := e.rConnector.GetClient().Set(
    ctx,
    getAlertKey(alert.AlertId, channel),
    getAlertValue(alert),
    duration,
  ).Result()
  return err
}