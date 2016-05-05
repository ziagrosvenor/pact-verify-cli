import Promise from "bluebird"
import {pipeP, tap} from "ramda"
import {bufferToStr} from "./utils"
import {readPactJsonFile, parsePactJsonContract} from "./get-pact-file"
import { createPactHelperFromContract,
 writePactHelperFile} from "./create-pact-helper"

export default function () {
  return pipeP(
    readPactJsonFile,
    bufferToStr,
    parsePactJsonContract,
    createPactHelperFromContract,
    writePactHelperFile
  )()
}
