package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/otelcol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to the configuration file")
)

func main() {
	flag.Parse()

	// Create logger with stdout output
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Create settings
	settings := otelcol.CollectorSettings{
		Factories: func() (otelcol.Factories, error) {
			return components, nil
		},
		BuildInfo: component.NewDefaultBuildInfo(),
		ConfigProviderSettings: otelcol.ConfigProviderSettings{
			ResolverSettings: confmap.ResolverSettings{
				URIs: []string{*configPath},
				ProviderFactories: []confmap.ProviderFactory{
					fileprovider.NewFactory(),
				},
			},
		},
		LoggingOptions: []zap.Option{
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
			}),
		},
	}

	// Create collector
	collector, err := otelcol.NewCollector(settings)
	if err != nil {
		logger.Fatal("Failed to create collector", zap.Error(err))
	}

	// Create context that will be canceled on SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal: %s\n", sig.String())
		logger.Info("Received signal", zap.String("signal", sig.String()))
		cancel()
	}()

	// Start the collector
	if err := collector.Run(ctx); err != nil {
		logger.Fatal("Collector run failed", zap.Error(err))
	}
}
