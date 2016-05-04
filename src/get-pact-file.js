import Promise from "bluebird"
const fs = Promise.promisifyAll(require("fs"))

export const readFileFromArgV = () => fs.readFileAsync(__dirname + "/../" +  process.argv[2])
export const bufferToStr =  (buffer) => buffer.toString()

export function parsePactJson(pactJson) {
  pactJson = pactJson.replace(/(.*),.*\n.*"s":\s".*\n/g, "$1");
  return jsonParse(pactJson)
}

function jsonParse(str) {
  try {
    return JSON.parse(str)
  } catch(err) {
    throw new Error("Failed to parse pact file" + err.toString())
  }
}
