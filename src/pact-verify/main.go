package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func buildPactHelperFromPactJson(pactFilePath string) string {
	var dir = getPwd()

	dat, err := ioutil.ReadFile(path.Join(dir, pactFilePath))
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
	err := ioutil.WriteFile(path.Join(ROOT_DIR, "/tmp/pact_helper.rb"), d1, 0644)
	check(err)
}

func printHelp() {
	cmd := exec.Command("pact-verify", "help")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
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
			Usage:       "Read a Pact file from `PATH` and process it",
			Destination: &pactFilePath,
		},

		cli.StringFlag{
			Name:        "provider, prov",
			Usage:       "The `URL` of the provider service to verify the pact with",
			Destination: &providerUrl,
		},

		cli.StringFlag{
			Name:        "setup, s",
			Usage:       "The `URL` of the provider state server - This is used to process provider states",
			Destination: &stateServerUrl,
		},
	}

	app.Name = "verify"
	app.Usage = "Command line interface for Pact verification"
	app.Action = func(c *cli.Context) error {

		if pactFilePath == "" {
			return cli.NewExitError("\nEXITED \nA Pact file path is required i.e /tmp/pacts/pact.json \n", 86)
		}

		if providerUrl == "" {
			return cli.NewExitError("\nEXITED \nA provider service URL is required\n", 86)
		}

		if stateServerUrl == "" {
			return cli.NewExitError("\nEXITED \nA provider states setup service URL is required\n", 86)
		}

		var ROOT_DIR = getSrcDir()
		var PWD = getPwd()

		var pactHelperStr = buildPactHelperFromPactJson(pactFilePath)
		writePactHelperFile(ROOT_DIR, pactHelperStr)

		cmd := exec.Command(
			"sh",
			path.Join(ROOT_DIR, "/bin/run-pact-verify.sh"),
			path.Join(PWD, pactFilePath),
			providerUrl,
			stateServerUrl,
			ROOT_DIR,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return nil
	}

	app.Run(os.Args)
}

var errorExitMessage = `
EXITING due to an error.

Please ensure all CLI options are correctly configured.
Try pact-verify help for more information.

ERROR:
`

func check(e error) {
	if e != nil {
		fmt.Println(errorExitMessage)
		log.Fatal(e)
	}
}

func getPwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	check(err)

	return dir
}

func getSrcDir() string {
	var env = os.Getenv("CLI_SRC_DIR")

	if env != "" {
		return env
	}

	return os.Getenv("GOPATH")
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
