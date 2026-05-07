package cmd

import (
	"encoding/json"
	"os"

	"github.com/mshindle/datastream"
	rl "github.com/mshindle/datastream/logger"
	"github.com/mshindle/mockdata/internal/chance"
	"github.com/spf13/cobra"
)

var crapsCmd = &cobra.Command{
	Use:   "craps",
	Short: "shoot craps",
	Long: `
Emulates a dice throw similar to the casino game Craps. Rolls 2 six-sided die.`,
	RunE: runCraps,
}

func init() {
	rootCmd.AddCommand(crapsCmd)
}

func runCraps(cmd *cobra.Command, args []string) error {
	c := chance.NewCup(chance.Die{Sided: 6}, chance.Die{Sided: 6})
	gen := datastream.GeneratorFunc[*chance.ThrowResult](func() *chance.ThrowResult { return c.Throw() })
	marsh := datastream.MarshalFunc[*chance.ThrowResult](func(tr *chance.ThrowResult) ([]byte, error) { return json.Marshal(tr) })
	sink := rl.Println(os.Stdout)
	opt := datastream.WithRateLimit[*chance.ThrowResult](5.0, 1)
	e := datastream.NewEngine[*chance.ThrowResult](gen, marsh, sink, opt)
	return runEngine(cmd.Context(), e)
}
