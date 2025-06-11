package main

import (
	"context"
	"log"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/debugexporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

var components = otelcol.Factories{
	Receivers: map[component.Type]receiver.Factory{
		otlpreceiver.NewFactory().Type(): otlpreceiver.NewFactory(),
	},
	Processors: map[component.Type]processor.Factory{
		component.MustNewType("metric_processor"): processor.NewFactory(
			component.MustNewType("metric_processor"),
			createMetricProcessorConfig,
			processor.WithMetrics(createMetricProcessor, component.StabilityLevelStable),
		),
		component.MustNewType("batch"): batchprocessor.NewFactory(),
	},
	Exporters: map[component.Type]exporter.Factory{
		debugexporter.NewFactory().Type(): debugexporter.NewFactory(),
		kafkaexporter.NewFactory().Type(): kafkaexporter.NewFactory(),
	},
}

func createMetricProcessorConfig() component.Config {
	return &struct {
		RedisAddr string `mapstructure:"redis_addr"`
	}{}
}

func createMetricProcessor(
	_ context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (processor.Metrics, error) {
	config := cfg.(*struct {
		RedisAddr string `mapstructure:"redis_addr"`
	})
	log.Println("config", config)
	return NewMetricProcessor(nextConsumer, set.Logger, config.RedisAddr), nil
}
