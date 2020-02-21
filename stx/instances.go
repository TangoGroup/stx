package stx

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/cue/parser"
	"github.com/logrusorgru/aurora"
)

type instanceHandler func(*build.Instance, *cue.Instance, cue.Value)

// GetBuildInstances loads and parses cue files and returns a list of build instances
func GetBuildInstances(args []string, pkg string) []*build.Instance {
	const syntaxVersion = -1000 + 13

	config := load.Config{
		Package: pkg,
		Context: build.NewContext(
			build.ParseFile(func(name string, src interface{}) (*ast.File, error) {
				return parser.ParseFile(name, src,
					parser.FromVersion(syntaxVersion),
					parser.ParseComments,
				)
			})),
	}
	if len(args) < 1 {
		args = append(args, "./...")
	}
	// load finds files based on args and passes those to build
	// buildInstances is a list of build.Instances, each has been parsed
	buildInstances := load.Instances(args, &config)
	return buildInstances
}

// Process iterates over instances applying the handler function for each
func Process(buildInstances []*build.Instance, exclude string, handler instanceHandler) {

	au := aurora.NewAurora(true) // TODO move to logger

	var excludeRegexp *regexp.Regexp
	var excludeRegexpErr error
	var wg sync.WaitGroup
	var binsts []*build.Instance

	if exclude != "" {
		excludeRegexp, excludeRegexpErr = regexp.Compile(exclude)
		if excludeRegexpErr != nil {
			fmt.Println(au.Red(excludeRegexpErr.Error()))
			os.Exit(1)
		}
	}

	if excludeRegexp != nil {
		// filter build instances
		for _, buildInstance := range buildInstances {
			if excludeRegexp.MatchString(buildInstance.DisplayPath) {
				continue
			}
			binsts = append(binsts, buildInstance)
		}
	} else {
		binsts = buildInstances
	}

	cueInstances := cue.Build(binsts)

	for i := 0; i < len(binsts); i++ {
		wg.Add(1)
		go func(buildInstance *build.Instance, cueInstance *cue.Instance) {
			// A cue instance defines a single configuration based on a collection of underlying CUE files.
			// cue.Build is designed to produce a single cue.Instance from n build.Instances
			// doing so however, would loose the connection between a stack and the build instance that
			// contains relevant path/file information related to the stack
			// here we cue.Build one at a time so we can maintain a 1:1:1:1 between
			// build.Instance, cue.Instance, cue.Value, and Stack

			if cueInstance.Err != nil {
				// parse errors will be exposed here
				// fmt.Println(au.Cyan(binst.DisplayPath))
				fmt.Println(cueInstance.Err, cueInstance.Err.Position())
				wg.Done()
				return
			}
			handler(buildInstance, cueInstance, cueInstance.Value())
			wg.Done()
		}(binsts[i], cueInstances[i])
	}

	wg.Wait()
}
