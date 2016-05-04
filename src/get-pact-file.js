import Promise from "bluebird"
const fs = Promise.promisifyAll(require("fs"))
import {jsonParse} from "./utils"

export const readPactJsonFile = () => fs.readFileAsync(__dirname + "/../" +  process.argv[2])

export function parsePactJsonContract(pactJson) {
  pactJson = pactJson.replace(/(.*),.*\n.*"s":\s".*\n/g, "$1");
  return jsonParse(pactJson)
}

