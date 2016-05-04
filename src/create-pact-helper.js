import Promise from "bluebird"
import {map, path} from "ramda"
import providerStateTemplate from "./provider-state-template"
import {bufferToStr} from "./get-pact-file"
const fs = Promise.promisifyAll(require("fs"))

export function writePactFile(fileStr) {
  return fs.writeFileAsync(__dirname + "/../tmp/pact_helper.rb", fileStr)
}

export function readPactHelperCore() {
  return fs.readFileAsync(__dirname + "/pact-helper-core.rb")
}

export function createPactHelperFromContract(contractObj) {
  return readPactHelperCore()
    .then(bufferToStr)
    .then((str) => {
      const providerStatesRubyStr = map(
        providerStateTemplate(contractObj.consumer.name),
        map(path(["provider_state"]), contractObj.interactions)
      )

      return str + "\n\n" + providerStatesRubyStr.join("\n")
    })
}
