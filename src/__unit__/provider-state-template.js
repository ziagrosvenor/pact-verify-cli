import {expect} from "chai"
import providerStateTemplate from "../provider-state-template"
import providerStateTemplateOutput from "./fixtures/provider-state-template-output"

const CONSUMER_NAME = "CONSUMER_NAME"
const providerStates = [
  "STATE_1",
  "STATE_2",
]

describe("providerStateTemplate", function() {
  it("Template output is correct", () => {
    const result = providerStateTemplate(CONSUMER_NAME, providerStates)
    expect(result)
      .to
      .equal(providerStateTemplateOutput)
  })
})
