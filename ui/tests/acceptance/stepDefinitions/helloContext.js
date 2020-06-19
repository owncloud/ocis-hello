const assert = require('assert')
const { client } = require('nightwatch-api')
const { Given, When, Then } = require('cucumber')

Given('the user browses to the hello page', function () {
  return client.page.helloPage().navigateAndWaitTillLoaded()
})

When('the user submits {string} to the Greet input', function (input) {
  return client.page.helloPage().inputName(input)
})

Then('{string} should be shown in the hello screen', async function (result) {
  const hello = await client.page.helloPage().getHelloOutput()
  assert.strictEqual(hello, result, 'The output on hello screen doesnt matches to ' + result)
})
