import {curry, map} from "ramda"

export const providerStateTemplate = curry(function providerStateTemplate(consumerName, providerState) {
  return `
  provider_state "${providerState}" do
    set_up do
      set_up_state "${consumerName}", "${providerState}"
    end
  end
  `
})

export default curry(function providerStatesTemplate(
    consumerName,
    providerStates
) {
  return  `
Pact.provider_states_for "${consumerName}" do
  ${map(providerStateTemplate(consumerName), providerStates)}
end
`
})

