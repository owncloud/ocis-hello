const { Given, When, Then} = require('@cucumber/cucumber')
// import expect for assertion
const { expect } = require('@playwright/test')
const { LoginPage } = require('../pageobjects/loginPage')
const { HelloPage } = require('../pageobjects/helloPage')

Given('user {string} has logged in using the webUI', async function (username) {
  const loginPage = new LoginPage(global.page)
  await loginPage.goto()
  await loginPage.login(username, username)
})

Given('the user browses to the hello page', async function () {
  const helloPage = new HelloPage(global.page)
  await helloPage.goto()
})

When('the user submits {string} to the Greet input', async function (input) {
  const helloPage = new HelloPage(global.page)
  await helloPage.inputName(input)
})

Then('{string} should be shown in the hello screen', async function (expectedResult) {
  const helloPage = new HelloPage(global.page)
  const result = await helloPage.getHelloOutput()
  expect(result).toBe(expectedResult)
})
