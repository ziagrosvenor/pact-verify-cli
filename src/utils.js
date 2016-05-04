export function jsonParse(str) {
  try {
    return JSON.parse(str)
  } catch(err) {
    throw new Error("Failed to parse JSON " + err.toString())
  }
}

export const bufferToStr =  (buffer) => buffer.toString()
