package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"github.com/TangoGroup/stx/stx"
	"github.com/ghodss/yaml"
)

func saveStackAsYml(stackName string, stack stx.Stack, buildInstance *build.Instance, cueValue cue.Value) string {
	dir := filepath.Clean(config.CueRoot + "/" + config.Xpt.YmlPath + "/" + stack.Profile)
	os.MkdirAll(dir, 0766)
	//fmt.Println(err)
	fileName := dir + "/" + stackName + ".cfn.yml"
	fmt.Printf("%s %s %s %s\n", au.White("Saving"), au.Magenta(stackName), au.White("⤏"), fileName)
	template := cueValue.Lookup("Stacks", stackName, "Template")
	yml, _ := yaml.Marshal(template)
	//fmt.Printf("YAML: %+v\n", string(yml))
	ioutil.WriteFile(fileName, yml, 0766)
	return fileName
}
