import Promise from "bluebird"
import {pipeP, tap} from "ramda"
import {readFileFromArgV, bufferToStr, parsePactJson} from "./get-pact-file"
import { createPactHelperFromContract,
 writePactFile} from "./create-pact-helper"

pipeP(
  readFileFromArgV,
  bufferToStr,
  parsePactJson,
  createPactHelperFromContract,
  writePactFile
)()
.catch((err) => {
  console.log(err.stack)
})
