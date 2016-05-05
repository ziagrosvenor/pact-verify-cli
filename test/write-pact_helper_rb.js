import Promise from "bluebird"
import {expect} from "chai"
import createPactHelperFromContract from "../src/create-pact-helper-from-contract"
import {bufferToStr} from "../src/utils"

const fs = Promise.promisifyAll(require("fs"))
const pactHelperFileFixure =
  bufferToStr(fs.readFileSync(__dirname + "/fixtures/pact_helper.rb"))

function readOutputPactHelper() {
  return fs.readFileAsync(__dirname + "/../tmp/pact_helper.rb")
    .then(bufferToStr)
}

describe("writing a pact_helper.rb based on pact-file.json", function() {
  before(() => {
    process.argv[2] = "test/fixtures/pact-file.json"
  })

  it("writes the file ../tmp/pact_helper.rb", () => {
    return createPactHelperFromContract()
      .then(() => readOutputPactHelper())
      .then((outputPactHelperStr) => {
        expect(outputPactHelperStr)
          .to
          .equal(pactHelperFileFixure)
      })
  })
})
