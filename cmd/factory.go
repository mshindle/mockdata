package cmd

import (
	"context"
	"fmt"
	stdlog "log"
	"os"

	"github.com/mshindle/datagen"
	"github.com/mshindle/datagen/elastic"
	"github.com/mshindle/datagen/kafka"
	"github.com/mshindle/datagen/logger"
	"github.com/mshindle/zlg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type publisherFactory func(ctx context.Context, config appConfig) (datagen.Publisher, error)

var factoryLookup = map[string]publisherFactory{
	"kafka":     initKafka,
	"elastic":   initElastic,
	"cloudlogs": initCloudlogs,
	"logs":      initLogs,
}

func initPublisher(ctx context.Context, cfg appConfig) (datagen.Publisher, error) {
	factory, ok := factoryLookup[cfg.Sink]
	if !ok {
		return nil, fmt.Errorf("invalid sink option: %s", cfg.Sink)
	}
	return factory(ctx, cfg)
}

func initKafka(ctx context.Context, cfg appConfig) (datagen.Publisher, error) {
	p, err := kafka.New(cfg.Kafka)
	if err != nil {
		log.Error().Msg("unable to create kafka client")
		return nil, err
	}

	// just count the number of successful messages published...
	go func() {
		var n int
		for _ = range p.Successes() {
			n++
			if n%10000 == 0 {
				log.Debug().Int("messages", n).Msg("messages sent")
			}
		}
	}()

	// log any errors talking to kafka
	go func() {
		for err := range p.Errors() {
			log.Error().Err(err).Msg("failed to deliver message")
		}
	}()

	return p, nil
}

func initElastic(ctx context.Context, cfg appConfig) (datagen.Publisher, error) {
	p, err := elastic.New(ctx, cfg.Elastic)
	if err != nil {
		log.Error().Msg("unable to create elastic client")
		return nil, err
	}
	err = p.ListIndices()
	if err != nil {
		log.Error().Msg("could not pull index alias")
		return nil, err
	}
	return p, nil
}

func initCloudlogs(ctx context.Context, cfg appConfig) (datagen.Publisher, error) {
	zlog, err := zlg.NewWriter(ctx, cfg.CloudLog.Parent, cfg.CloudLog.LogID)
	if err != nil {
		return nil, err
	}
	// TODO: need to properly close zlg

	std := stdlog.Default()
	std.SetFlags(0)
	std.SetOutput(zlog)
	return logger.LoggerPublisher{Logger: std}, nil
}

func initLogs(ctx context.Context, cfg appConfig) (datagen.Publisher, error) {
	zlog := zerolog.New(os.Stdout).With().Timestamp().Str("event", "publish").Logger()
	std := stdlog.Default()
	std.SetFlags(0)
	std.SetOutput(zlog)

	return logger.LoggerPublisher{Logger: std}, nil
}
