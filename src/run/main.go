package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func PactVerify(rootDir string, pwd string, pactFilePath string, providerUrl string, stateServerUrl string) {
	cmd := exec.Command(
		"sh",
		path.Join(rootDir, "/bin/run-pact-verify.sh"),
		path.Join(pwd, pactFilePath),
		providerUrl,
		stateServerUrl,
		rootDir,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func BuildPactHelperFromPactJson(dir string, pactFilePath string) string {
	dat, err := ioutil.ReadFile(path.Join(dir, pactFilePath))
	check(err)
	var pact map[string]interface{}

	if err := json.Unmarshal(dat, &pact); err != nil {
		panic(err)
	}

	var consumer = pact["consumer"].(map[string]interface{})
	var consumerName = consumer["name"].(string)
	var interactions = pact["interactions"].([]interface{})

	fmt.Println(interactions)

	return TemplatePactHelper(consumerName, interactions)
}

func WritePactHelperFile(ROOT_DIR string, pactHelperStr string) {
	d1 := []byte(pactHelperStr)
	err := ioutil.WriteFile(path.Join(ROOT_DIR, "/tmp/pact_helper.rb"), d1, 0644)
	check(err)
}

func TemplatePactHelper(consumerName string, interactions []interface{}) string {
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

func check(e error) {
	if e != nil {
		fmt.Println(errorExitMessage)
		log.Fatal(e)
	}
}

var errorExitMessage = `
EXITING due to an error.

Please ensure all CLI options are correctly configured.
Try pact-verify help for more information.

ERROR:
`

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

var pactHelperBaseStr = `require 'faraday'
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
