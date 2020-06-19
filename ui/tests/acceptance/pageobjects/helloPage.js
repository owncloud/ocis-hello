module.exports = {
  url: function () {
    return this.api.launchUrl + '/#/hello'
  },

  commands: {
    navigateAndWaitTillLoaded: async function () {
      const url = this.url()
      return this.navigate(url).waitForElementVisible('@accountInput')
    },

    inputName: async function (name) {
      await this
        .waitForElementVisible('@accountInput')
        .clearValue('@accountInput')
        .setValue('@accountInput', name)
      return this.waitForElementVisible('@submitButton')
        .click('@submitButton')
    },

    getHelloOutput: async function () {
      let output
      await this.waitForElementVisible('@helloText')
        .getText('@helloText', (result) => {
          output = result
        })
      return output.value
    }
  },

  elements: {
    accountInput: {
      selector: "//input[@placeholder='Your name']",
      locateStrategy: 'xpath'
    },
    submitButton: {
      selector: "//button[contains(@class, 'uk-button-primary')]",
      locateStrategy: 'xpath'
    },
    helloText: {
      selector: "//*[@class='uk-text-lead']",
      locateStrategy: 'xpath'
    }
  }
}
