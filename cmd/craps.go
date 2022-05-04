package cmd

import (
	"github.com/mshindle/datagen"
	"github.com/spf13/cobra"

	"github.com/mshindle/mockdata/generators"
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
	c := generators.NewCup(6, 2)
	g := datagen.GeneratorFunc(func() datagen.Event {
		return c.Throw()
	})
	return executeEngine(g)
}
