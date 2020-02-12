package cmd

import (
	"fmt"

	"github.com/TangoGroup/stx/stx"
	"github.com/logrusorgru/aurora"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/pkg/encoding/yaml"
	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the Cue output as YAML",
	Long:  `yada yada yada`,
	Run: func(cmd *cobra.Command, args []string) {
		au := aurora.NewAurora(true)
		buildInstances := stx.GetBuildInstances(args, "cfn")
		stx.Process(buildInstances, func(buildInstance *build.Instance, cueInstance *cue.Instance, cueValue cue.Value) {
			fmt.Println(au.Cyan(buildInstance.DisplayPath))

			yml, ymlErr := yaml.Marshal(cueValue)

			if ymlErr != nil {
				fmt.Println(au.Red(ymlErr.Error()))
			} else {
				fmt.Printf("%s\n", string(yml))
			}

		})
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	// TODO add flag to skip/hide errors
	// TODO add flag to show only errors
}
