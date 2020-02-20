package cmd

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"github.com/TangoGroup/stx/stx"
	"github.com/spf13/cobra"
)

// xptCmd represents the xpt command
var xptCmd = &cobra.Command{
	Use:   "xpt",
	Short: "eXPorTs cue templates that implement the Stacks:[] pattern.",
	Long:  `Yada yada yada.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInstances := stx.GetBuildInstances(args, "cfn")
		stx.Process(buildInstances, flags.exclude, func(buildInstance *build.Instance, cueInstance *cue.Instance, cueValue cue.Value) {
			stacks := stx.GetStacks(cueValue)
			if stacks != nil {
				for stackName, stack := range stacks {
					saveStackAsYml(stackName, stack, buildInstance, cueValue)
				}
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(xptCmd)
}
