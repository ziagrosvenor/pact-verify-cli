import Promise from "bluebird"
import {map, path} from "ramda"
import providerStateTemplate from "./provider-state-template"
import {bufferToStr} from "./utils"
const fs = Promise.promisifyAll(require("fs"))

export function writePactHelperFile(fileStr) {
  return fs.writeFileAsync(__dirname + "/../tmp/pact_helper.rb", fileStr)
}

export function readPactHelperBase() {
  return fs.readFileAsync(__dirname + "/pact_helper_base.rb")
}

export function createPactHelperFromContract(contractObj) {
  return readPactHelperBase()
    .then(bufferToStr)
    .then((str) => {
      const providerStatesRubyStr = map(
        providerStateTemplate(contractObj.consumer.name),
        map(path(["provider_state"]), contractObj.interactions)
      )

      return str + "\n\n" + providerStatesRubyStr.join("\n")
    })
}
