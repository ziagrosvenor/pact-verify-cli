import Promise from "bluebird"
import {pipeP, tap} from "ramda"
import {bufferToStr} from "./utils"
import {readPactJsonFile, parsePactJsonContract} from "./get-pact-file"
import { createPactHelperFromContract,
 writePactHelperFile} from "./create-pact-helper"

pipeP(
  readPactJsonFile,
  bufferToStr,
  parsePactJsonContract,
  createPactHelperFromContract,
  writePactHelperFile
)()
.catch((err) => {
  console.log(err.stack)
})
