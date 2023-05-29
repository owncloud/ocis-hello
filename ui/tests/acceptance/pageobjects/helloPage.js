const appSwitcherButtonSelector = '#_appSwitcherButton'
const appSwitcherDropdownSelector = '#app-switcher-dropdown'
const appName = 'Hello'
const accountInputSelector = "//input[@placeholder='Your name']"
const submitButtonSelector = "//button[@type='submit']"
const helloTextSelector = "//*[@class='uk-text-lead']"

exports.HelloPage = class HelloPage {
  constructor (page) {
    this.page = page
  }

  async goto () {
    await this.page.goto('/')
    await this.page.locator(appSwitcherButtonSelector).click()
    await this.page.locator(appSwitcherDropdownSelector).getByText(appName).click()
  }

  async inputName (name) {
    await this.page.locator(accountInputSelector).fill(name)
    await this.page.locator(submitButtonSelector).click()
  }

  async getHelloOutput () {
    return await this.page.locator(helloTextSelector).innerText()
  }
}
