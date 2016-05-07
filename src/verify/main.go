package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var pactHelperTemplate = `
%[1]v

Pact.provider_states_for "%[2]v" do
  %[3]v
end
`

var template = `
  provider_state "%[1]v" do
    set_up do
      set_up_state "%[2]v", "%[1]v"
    end
  end
`

func getRootDirPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	check(err)

	return dir
}

func buildPactHelperFromPactJson(pactFilePath string) string {
	var dir = getRootDirPath()

	dat, err := ioutil.ReadFile(dir + pactFilePath)
	check(err)

	pactHelperBase, err := ioutil.ReadFile(dir + "/src/pact_helper_base.rb")
	check(err)

	pactHelperBaseStr := string(pactHelperBase)

	var pact map[string]interface{}

	if err := json.Unmarshal(dat, &pact); err != nil {
		panic(err)
	}

	var consumer = pact["consumer"].(map[string]interface{})
	var consumerName = consumer["name"].(string)

	var interactions = pact["interactions"].([]interface{})

	var setupStateMethodCalls bytes.Buffer

	for _, element := range interactions {
		var interaction = element.(map[string]interface{})
		var providerState = interaction["provider_state"].(string)

		setupStateMethodCalls.WriteString(fmt.Sprintf(template, providerState, consumerName))
	}

	var pactHelperRubyStr = fmt.Sprintf(
		pactHelperTemplate,
		pactHelperBaseStr,
		consumerName,
		setupStateMethodCalls.String(),
	)

	return pactHelperRubyStr
}

func writePactHelperFile(pactHelperStr string) {
	var dir = getRootDirPath()
	d1 := []byte(pactHelperStr)
	err := ioutil.WriteFile(dir+"/tmp/pact_helper.rb", d1, 0644)
	check(err)
}

func main() {
	var pactFilePath string
	var providerUrl string
	var stateServerUrl string

	app := cli.NewApp()
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "pact",
			Usage:       "Load Pact JSON from `FILE`",
			Destination: &pactFilePath,
		},

		cli.StringFlag{
			Name:        "provider, prov",
			Usage:       "The URL of the provider service that the pact will be verified against",
			Destination: &providerUrl,
		},

		cli.StringFlag{
			Name:        "setup, s",
			Usage:       "The URL of the provider state server - This is used to process provider states",
			Destination: &stateServerUrl,
		},
	}

	app.Name = "verify"
	app.Usage = "Command line interface for Pact verification"
	app.Action = func(c *cli.Context) error {

		if pactFilePath == "" {
			fmt.Printf("\nEXITED \nPact file path required i.e /tmp/pacts/pact.json \n\n")
			return nil
		}

		if providerUrl == "" {
			fmt.Printf("\nEXITED \nProvider url required\n\n")
			return nil
		}

		var dir = getRootDirPath()
		var pactHelperStr = buildPactHelperFromPactJson(pactFilePath)
		writePactHelperFile(pactHelperStr)
		cmd := exec.Command("sh", dir+"/bin/run-pact-verify.sh", "."+pactFilePath, providerUrl, stateServerUrl)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return nil
	}

	app.Run(os.Args)
}
