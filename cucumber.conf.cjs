const { Before, BeforeAll, AfterAll, After, setDefaultTimeout } = require('@cucumber/cucumber');
const { chromium } = require('@playwright/test')

setDefaultTimeout(60000)

// launch the browser
BeforeAll(async function () {
  global.browser = await chromium.launch({
    headless: false,
    slowMo: 1000 // TODO set to 0 in CI
  })
})

// close the browser
AfterAll(async function () {
  await global.browser.close()
})

// Create a new browser context and page per scenario
Before(async function () {
  global.context = await global.browser.newContext(
    {
      ignoreHTTPSErrors: true,
      baseURL: 'https://host.docker.internal:9200' // TODO make flexible
    }
  )
  global.page = await global.context.newPage()
})

// Cleanup after each scenario
After(async function () {
  await global.page.close()
  await global.context.close()
})
