const usernameSelector = '#oc-login-username'
const passwordSelector = '#oc-login-password'
const buttonSelector = 'button[type="submit"]'

exports.LoginPage = class LoginPage {
  constructor (page) {
    this.page = page
  }

  async goto () {
    await this.page.goto('/')
  }

  async login (username, password) {
    await this.page.locator(usernameSelector).fill(username)
    await this.page.locator(passwordSelector).fill(password)
    await this.page.locator(buttonSelector).click()
  }
}
