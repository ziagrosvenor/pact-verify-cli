package run

import (
	"log"
	"os"
	"regexp"
	"testing"
)

var srcDir = os.Getenv("CLI_SRC_DIR")

func TestBuildPactHelperFromPactJson(t *testing.T) {
	var pactStr string
	pactStr = BuildPactHelperFromPactJson(srcDir, "./test/fixtures/pact-file.json")
	matched, err := regexp.MatchString(pactStr, outputPactHelper)

	if err != nil {
		log.Fatal(err)
	}

	if !matched {
		t.Error("Expected:\n"+outputPactHelper+"\n\n", pactStr)
	}
}

var outputPactHelper = `
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


Pact.provider_states_for "CONSUMER" do
  
  provider_state "active_user" do
    set_up do
      set_up_state "CONSUMER", "active_user"
    end
  end

end`
