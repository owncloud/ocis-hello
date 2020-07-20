'use strict'

const { pactWith } = require('jest-pact')
// eslint-disable-next-line camelcase
const { Hello_Greet } = require('../../ui/client/hello/index')

pactWith({ consumer: 'OCIS-Hello_UI', provider: 'OCIS-Hello_API' }, provider => {
  describe('Hello API', () => {
    const RESPONSE_DATA = {
      message: 'Hello Milan'
    }

    const REQUEST_DATA = {
      name: 'Milan'
    }
    const successResponse = {
      status: 201,
      headers: {
        'Content-Type': 'application/json; charset=utf-8'
      },
      body: RESPONSE_DATA
    }

    const greetRequest = {
      uponReceiving: 'a greeting',
      withRequest: {
        method: 'POST',
        path: '/api/v0/greet',
        body: REQUEST_DATA
      }
    }

    beforeEach(() => {
      const interaction = {
        state: 'i am greeted with my name',
        ...greetRequest,
        willRespondWith: successResponse
      }
      return provider.addInteraction(interaction)
    })

    // add expectations
    it('returns a successful body', () => {
      return Hello_Greet({
        $domain: provider.mockService.baseUrl,
        body: REQUEST_DATA
      }).then(response => {
        expect(response.data).toEqual(RESPONSE_DATA)
      })
    })
  })
})
