const { Before, BeforeAll, AfterAll, After, setDefaultTimeout } = require('@cucumber/cucumber');
const { chromium } = require('@playwright/test')

setDefaultTimeout(60000)

// launch the browser
BeforeAll(async function () {
  global.browser = await chromium.launch({
    headless: process.env.CI === 'true',
    slowMo: process.env.CI === 'true' ? 0 : 1000
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
      baseURL: process.env.SERVER_HOST ?? 'https://localhost:9200'
    }
  )
  global.page = await global.context.newPage()
})

// Cleanup after each scenario
After(async function () {
  await global.page.close()
  await global.context.close()
})
