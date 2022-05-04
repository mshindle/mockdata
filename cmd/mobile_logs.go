package cmd

import (
	"github.com/mshindle/datagen"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mshindle/mockdata/generators"
)

// mobileLogsCmd represents the kafka command
var mobileLogsCmd = &cobra.Command{
	Use:   "mobilelogs",
	Short: "generate mobile log data",
	RunE:  runMobileLogs,
}

func init() {
	rootCmd.AddCommand(mobileLogsCmd)
}

func runMobileLogs(cmd *cobra.Command, args []string) error {
	generator := datagen.GeneratorFunc(
		func() datagen.Event {
			e, err := generators.MockMobileLog()
			if err != nil {
				log.Error().Err(err).Msg("unable to generate fake data")
				return nil
			}
			return e
		},
	)
	return executeEngine(generator)
}
