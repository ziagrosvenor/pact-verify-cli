# Pact Verify CLI

## Features / Objectives

- Retrieve Pact JSON contracts from the file system.

- Handle necessary Pact configuration to verify the contracts against the provider service.

- Verify Pact JSON contracts against the provider service.

- Delegate responsibility for handling the provider states over to a provider state server. This could be located inside the repository for the provider service.

## NPM Installation
    $ gem install bundler
    $ npm i pact-verify-cli -g

      > pact-verify-cli@1.3.3 postinstall /Users/zia.grosvenor/.nvm/versions/node/v0.12.7/lib/node_modules/pact-verify-cli
      > bundle install && mkdir ./tmp

    $ export GOPATH=/Users/zia.grosvenor/.nvm/versions/node/v0.12.7/lib/node_modules/pact-verify-cli

## Installation via Github
    $ git clone git@github.com:ziagrosvenor/pact-verify-cli.git
    $ cd ./pact-verify-cli
    $ bundle install
    $ export GOPATH=`pwd`
    $ PATH=$PATH:$GOPATH/bin

You may want to add this to your ~/.bash_profile
    $ export GOPATH=/path/to/pact-verify-cli
    $ PATH=$PATH:$GOPATH/bin

## CLI

#### Verify a Pact file with the provider service
The `provider_state` and `consumer` are sent to the setup URL in a POST request.

    $ pact-verify --pact /tmp/pacts/pact-file.json --provider http://localhost:3000 --setup http://localhost:3001

#### Help

    $ pact-verify help
    NAME:
        verify - Command line interface for Pact verification

    USAGE:
        verify [global options] command [command options] [arguments...]

    VERSION:
        0.0.0

    COMMANDS:
    GLOBAL OPTIONS:
        --pact FILE				Load Pact JSON from FILE
        --provider value, --prov value	The URL of the provider service that the pact will be verified against
        --setup value, -s value		The URL of the provider state server - This is used to process provider states
        --help, -h				show help
        --version, -v			print the version


This is interface was designed in the documentation for [the `provider-state-server-poc` repository of Github](https://github.com/bethesque/provider-state-server-poc)

### Terminology
[Pact terminology wiki](https://github.com/realestate-com-au/pact/wiki/Terminology)

### Provider state setup server
This will receive the consumers name and the required `provider_state` for the next interaction that `pact:verify` will test.

#### Example endpoint
```javascript
import server from "./server"
import setupUserInDatabase from "./db"

const providerStates = {
  CONSUMER_NAME: {
    ACTIVE_USER: () => {
      return setupUserInDatabase()
    }
  }
}

server.post("/setup", function postSetupState(req, res) {
  const {consumer, provider_state} = req.body

  providerStates[consumer][provider_state]
    .then(() => res.send(200))
    .catch(() => res.send(500))
})
```

#### Resources / Related Repositories
- [pact](https://github.com/realestate-com-au/pact)
- [pact-foundation](https://github.com/pact-foundation)
- [provider-state-server-poc](https://github.com/bethesque/provider-state-server-poc)
- [pact-consumer-js-dsl](https://github.com/DiUS/pact-consumer-js-dsl)
- [pact-node](https://github.com/pact-foundation/pact-node)
