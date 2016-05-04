import {curry} from "ramda"

export default curry(function providerStateTemplate(
    consumerName,
    providerState
) {
  return  `
Pact.provider_states_for "${consumerName}" do
  provider_state "${providerState}" do
    set_up do
      set_up_state "${consumerName}", "${providerState}"
    end
  end
end
`
})
