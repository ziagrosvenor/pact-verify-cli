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

var pactHelperBaseStr = `
require 'faraday'
require 'cgi'

PROVIDER_STATE_SERVER_SET_UP_URL = ENV["SETUP_SERVER_URL"] 

# Responsible for making the call to the provider state server to set up the state
module ProviderStateServerClient

  def set_up_state consumer_name, provider_state
    puts "Setting up provider state '#{provider_state}' using provider state server at #{PROVIDER_STATE_SERVER_SET_UP_URL}"
    Faraday.post(PROVIDER_STATE_SERVER_SET_UP_URL, {"consumer" => consumer_name, "provider_state" => provider_state })
  end

end

Pact.configure do | config |
  config.include ProviderStateServerClient
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

func writePactHelperFile(ROOT_DIR string, pactHelperStr string) {
	d1 := []byte(pactHelperStr)
	err := ioutil.WriteFile(ROOT_DIR+"/tmp/pact_helper.rb", d1, 0644)
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
			Usage:       "Read from disc and process a Pact JSON file from `PATH`",
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
			fmt.Printf("\nEXITED \nA Pact file path required i.e /tmp/pacts/pact.json \n\n")

			cmd := exec.Command("verify", "help")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			return nil
		}

		if providerUrl == "" {
			fmt.Printf("\nEXITED \nProvider url required\n\n")

			cmd := exec.Command("verify", "help")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			return nil
		}

		var ROOT_DIR = os.Getenv("GOPATH")
		var PWD = getRootDirPath()

		var pactHelperStr = buildPactHelperFromPactJson(pactFilePath)
		writePactHelperFile(ROOT_DIR, pactHelperStr)

		cmd := exec.Command(
			"sh",
			ROOT_DIR+"/bin/run-pact-verify.sh",
			PWD+"/"+pactFilePath,
			providerUrl,
			stateServerUrl,
			os.Getenv("GOPATH"),
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return nil
	}

	app.Run(os.Args)
}
