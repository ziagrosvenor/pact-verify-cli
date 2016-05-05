import {expect} from "chai"
import {omit} from "ramda"
import {parsePactJsonContract} from "../get-pact-file"

const mockPactJsonStr = `{
  "mockProp": "mockProp",
  "s": "will remove this"
}`

describe("parsePactJsonContract", function() {

  /** This is a fix for the regex value
   * of the 's' property causing a JSON parse error
   **/
  it("Remove the s property before parsing the JSON", () => {

    const result = parsePactJsonContract(mockPactJsonStr)
    expect(result)
      .to
      .eql(omit(["s"], JSON.parse(mockPactJsonStr)))
  })
})
