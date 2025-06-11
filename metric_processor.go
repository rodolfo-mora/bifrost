package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// MetricProcessor is a processor that stores metric names in Redis
type MetricProcessor struct {
	nextConsumer consumer.Metrics
	logger       *zap.Logger
	redisClient  *redis.Client
}

// NewMetricProcessor creates a new MetricProcessor
func NewMetricProcessor(nextConsumer consumer.Metrics, logger *zap.Logger, redisAddr string) *MetricProcessor {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &MetricProcessor{
		nextConsumer: nextConsumer,
		logger:       logger,
		redisClient:  client,
	}
}

func (p *MetricProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (p *MetricProcessor) Start(ctx context.Context, host component.Host) error {
	// Test Redis connection
	_, err := p.redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}

func (p *MetricProcessor) Shutdown(ctx context.Context) error {
	return p.redisClient.Close()
}

func (p *MetricProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	// Process each resource metrics
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rm := md.ResourceMetrics().At(i)

		// Process each scope metrics
		for j := 0; j < rm.ScopeMetrics().Len(); j++ {
			sm := rm.ScopeMetrics().At(j)

			// Process each metric
			for k := 0; k < sm.Metrics().Len(); k++ {
				metric := sm.Metrics().At(k)
				metricName := metric.Name()
				p.logger.Info(metricName)

				// Store metric name in Redis with timestamp
				key := fmt.Sprintf("metrics:%s", metricName)
				value := time.Now().Format(time.RFC3339)

				err := p.redisClient.Set(ctx, key, value, 24*time.Hour).Err()
				if err != nil {
					p.logger.Error("Failed to store metric in Redis",
						zap.String("metric", metricName),
						zap.Error(err))
				} else {
					p.logger.Info("Stored metric in Redis",
						zap.String("metric", metricName))
				}
			}
		}
	}

	// Pass the metrics to the next consumer
	return p.nextConsumer.ConsumeMetrics(ctx, md)
}
