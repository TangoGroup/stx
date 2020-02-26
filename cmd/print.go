package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/TangoGroup/stx/stx"

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

		if flags.PrintOnlyErrors && flags.PrintHideErrors {
			fmt.Println(au.Red("Cannot show only errors while hiding them."))
			os.Exit(1)
		}
		totalErrors := 0
		buildInstances := stx.GetBuildInstances(args, "cfn")
		stx.Process(buildInstances, flags, func(buildInstance *build.Instance, cueInstance *cue.Instance, cueValue cue.Value) {

			valueToMarshal := cueValue
			stacks := stx.GetStacks(cueValue, flags)

			for stackName := range stacks {
				var path []string
				var displayPath string
				if flags.PrintPath != "" {
					path = []string{"Stacks", stackName}
					path = append(path, strings.Split(flags.PrintPath, ".")...)
					valueToMarshal = cueValue.Lookup(path...)
					if valueToMarshal.Err() != nil {
						continue
					}
					displayPath = strings.Join(path, ".") + ":\n"
				}
				yml, ymlErr := yaml.Marshal(valueToMarshal)

				if ymlErr != nil {
					totalErrors++
					if !flags.PrintHideErrors {
						fmt.Println(au.Cyan(buildInstance.DisplayPath))
						fmt.Println(au.Red(ymlErr.Error()))
					}
				} else {
					if !flags.PrintOnlyErrors {
						fmt.Println(au.Cyan(buildInstance.DisplayPath))
						fmt.Printf("%s\n", displayPath+string(yml))
					}
				}
			}
		})

		if !flags.PrintHideErrors && totalErrors > 0 {
			fmt.Println("Total errors: ", totalErrors)
		}
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	// TODO add flag to skip/hide errors

	printCmd.Flags().BoolVar(&flags.PrintOnlyErrors, "only-errors", false, "Only print errors. Cannot be used in concjunction with --hide-errors")
	printCmd.Flags().BoolVar(&flags.PrintHideErrors, "hide-errors", false, "Hide errors. Cannot be used in concjunction with --only-errors")
	printCmd.Flags().StringVarP(&flags.PrintPath, "path", "p", "", "Dot-notation style path to key to print. Eg: Template.Resources.Alb or Template.Outputs")

}
