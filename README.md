# Pact Verify CLI

## Features / Objectives

- Retrieve Pact JSON contracts from the file system.

- Handle necessary Pact configuration to verify the contracts against the provider service.

- Verify Pact JSON contracts against the provider service.

- Delegate responsibility for handling the provider states over to a provider state server. This could be located inside the repository for the provider service.

### Installation
    bundle install
    npm i
    
### Terminology
[Pact terminology wiki](https://github.com/realestate-com-au/pact/wiki/Terminology)
    
### Current CLI implementation
    sh ./pact-verify ./tmp/pact-file.json http://localhost:3000 http://localhost:3001/setup

### Preferred CLI implementation
    pact-verify --pact=./tmp/pact-file.json
    --provider-url=http://localhost:3000
    --setup-url=http://localhost:3001/setup

This is interface was designed in the documentation for [the `provider-state-server-poc` repository of Github](https://github.com/bethesque/provider-state-server-poc)

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
  const {consumerName, providerState} = req.body
  
  providerStates[consumerName][providerState]
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
