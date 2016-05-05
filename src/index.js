import createPactHelperFromContract from "./create-pact-helper-from-contract"

createPactHelperFromContract()
  .catch((err) => {
    console.log(err.stack)
  })
