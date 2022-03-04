module.exports = {
  url: function () {
    return this.api.launchUrl + '/#/hello'
  },

  commands: {
    navigateAndWaitTillLoaded: async function () {
      const url = this.url()
      return this.navigate(url).waitForElementVisible('@helloInput')
    },

    inputName: async function (name) {
      await this
        .waitForElementVisible('@helloInput')
        .clearValue('@helloInput')
        .setValue('@helloInput', name)
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
    helloInput: {
      selector: "//input[@placeholder='Your name']",
      locateStrategy: 'xpath'
    },
    submitButton: {
      selector: "//button[contains(@class, 'oc-button-primary')]",
      locateStrategy: 'xpath'
    },
    helloText: {
      selector: "//*[@class='oc-text-lead']",
      locateStrategy: 'xpath'
    }
  }
}
