package cmd

import (
	"encoding/json"
	"os"

	"github.com/mshindle/mockdata/internal/systems"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mshindle/datastream"
	rl "github.com/mshindle/datastream/logger"
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
	// 1. Generator (Sequential & simple)
	gen := func() *systems.MobileLog {
		e, err := systems.MockMobileLog()
		if err != nil {
			log.Error().Err(err).Msg("generation failed")
		}
		return e
	}

	// 2. Marshaller (JSON)
	marsh := func(ml *systems.MobileLog) ([]byte, error) {
		return json.Marshal(ml)
	}

	engine := datastream.NewEngine[*systems.MobileLog](
		gen,
		marsh,
		rl.Println(os.Stdout),
		datastream.WithRateLimit[*systems.MobileLog](5.0, 1),
	)
	return runEngine(cmd.Context(), engine)
}
