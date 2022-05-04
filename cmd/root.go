package cmd

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/mshindle/datagen"
	"github.com/mshindle/datagen/elastic"
	"github.com/mshindle/datagen/kafka"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type appConfig struct {
	Debug      bool
	Generators int
	Publishers int
	Sink       string
	Kafka      kafka.Config
	Elastic    elastic.Config
	CloudLog   struct {
		Parent string
		LogID  string
	}
}

var cfgFile string
var cfgApp appConfig
var publisher datagen.Publisher

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "mockdata",
	Short:              "generate mock data and publish it to an endpoint",
	Long:               ``,
	PersistentPreRunE:  mockSetup,
	PersistentPostRunE: mockTearDown,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.datagen.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "sets log level to debug")
	rootCmd.PersistentFlags().IntP("generators", "g", 1, "set the number of generators")
	rootCmd.PersistentFlags().IntP("publishers", "p", 1, "set the number of publishers")
	rootCmd.PersistentFlags().String("sink", "log", "set the publish destination: logs, cloudlogs, kafka, or elastic")
	rootCmd.PersistentFlags().String("cloud_project", "", "specify the google project id to receive logs")
	rootCmd.PersistentFlags().String("name", "sample-log", "sets the name of the log to write to")

	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("generators", rootCmd.PersistentFlags().Lookup("generators"))
	_ = viper.BindPFlag("publishers", rootCmd.PersistentFlags().Lookup("publishers"))
	_ = viper.BindPFlag("sink", rootCmd.PersistentFlags().Lookup("sink"))
	_ = viper.BindPFlag("cloudlog.parent", rootCmd.PersistentFlags().Lookup("cloud_project"))
	_ = viper.BindPFlag("cloudlog.logID", rootCmd.PersistentFlags().Lookup("name"))

	viper.SetDefault("elastic.flushbytes", 5e+6)
	viper.SetDefault("elastic.flushinterval", 3*time.Second)
	viper.SetDefault("elastic.numworkers", runtime.NumCPU())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("exiting application...")
		}

		// Search config in home directory with name ".datagen" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yml")
		viper.SetConfigName(".mockdata")
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Str("file", viper.ConfigFileUsed()).Msg("using config file")
	}
}

func mockSetup(cmd *cobra.Command, args []string) error {
	err := viper.Unmarshal(&cfgApp)
	if err != nil {
		return err
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfgApp.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// create publisher
	publisher, err = initPublisher(cmd.Context(), cfgApp)
	return err
}

func mockTearDown(cmd *cobra.Command, args []string) error {
	p, ok := publisher.(*kafka.Service)
	if !ok {
		return nil
	}
	return p.Close()
}

// executeEngine executes the running of the engine and wrapping it around an
// os.Signal so the process can be killed cleanly from the cmdline
func executeEngine(generator datagen.Generator) error {
	engine, _ := datagen.NewEngine(
		generator,
		publisher,
		datagen.WithNumGenerators(cfgApp.Generators),
		datagen.WithNumPublishers(cfgApp.Publishers),
	)
	done, _ := engine.Run()
	defer close(done)

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down generation")

	// tell the engine to stop....
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	return nil
}
