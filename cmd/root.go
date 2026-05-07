package cmd

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mshindle/datastream"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel = zerolog.InfoLevel
	Version         = "0.1.0"
	serviceName     = "mockdata"
)

type appConfig struct {
	Log struct {
		Level   string `mapstructure:"level"`
		Console bool   `mapstructure:"console"`
	} `mapstructure:"log"`
	//Kafka   kafka.Config   `mapstructure:"kafka"`
	//Elastic elastic.Config `mapstructure:"elastic"`
}

var cfgFile string
var cfgApp appConfig
var v = viper.New()
var logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(defaultLogLevel)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "mockdata",
	Short:              "generate mock data and publish it to an endpoint",
	Long:               ``,
	PersistentPreRunE:  mockSetup,
	PersistentPostRunE: nil,
	SilenceErrors:      true,
	SilenceUsage:       true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("exiting application...")
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.console", false)
	v.SetDefault("elastic.flushbytes", 5e+6)
	v.SetDefault("elastic.flushinterval", 3*time.Second)
	v.SetDefault("elastic.numworkers", runtime.NumCPU())

	zerolog.TimeFieldFormat = time.RFC3339Nano
}

// initConfig reads in the config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		err := godotenv.Load(cfgFile)
		if err != nil {
			logger.Error().Err(err).Str("cfgFile", cfgFile).Msg("failed to load config file")
		}
	}
	err := godotenv.Load()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load ./.env file; skipping")
	}

	// 3. Setup Environment Variable Logic
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("MD") // Prepends "DRONE_" to all env lookups
	v.AutomaticEnv()
}

func mockSetup(cmd *cobra.Command, args []string) error {
	err := v.Unmarshal(&cfgApp)
	if err != nil {
		return err
	}

	// configure logging
	if cfgApp.Log.Console {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	lvl, err := zerolog.ParseLevel(cfgApp.Log.Level)
	if err == nil && lvl != logger.GetLevel() {
		logger = logger.Level(lvl)
	}
	logger.Info().Err(err).Str("min_level", logger.GetLevel().String()).Msg("minimum logging level")
	log.Logger = logger

	return nil
}

// executeEngine runs the datastream pipeline with a cancellation signal.
func runEngine[T any](ctx context.Context, engine *datastream.Engine[T]) error {
	// Create a context that reacts to interrupt signals
	notifyCtx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	log.Info().Msg("starting data stream...")
	return engine.Run(notifyCtx)
}
