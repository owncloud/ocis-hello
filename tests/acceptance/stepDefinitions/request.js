const { When, Then } = require('cucumber')
const fetch = require('node-fetch')
const assert = require('assert')
const { BASE_URL } = require('../setup.js')

let httpResponse

When('sending {string} to {string} with', function (method, path, table) {
  const requestData = table.rowsHash()
  let body
  let headers
  if (requestData.body !== undefined) {
    body = requestData.body
  }
  if (requestData.headers !== undefined) {
    headers = JSON.parse(requestData.headers)
  }
  return fetch(
    BASE_URL + path,
    { method: method, body: body, headers: headers }
  ).then(function (res) {
    httpResponse = res
  })
})

When('sending {string} to {string}', function (method, path) {
  return fetch(
    BASE_URL + path,
    { method: method }
  ).then(function (res) {
    httpResponse = res
  })
})

Then('response http status code should be {int}', function (statusCode) {
  return assert.strictEqual(statusCode, httpResponse.status)
})

Then('response body should be', async function (body) {
  return assert.strictEqual(body, await httpResponse.text())
})
