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

    const REQUEST_DATA_INVALID_FIELD = {
      surname: 'Milan'
    }

    const RESPONSE_DATA_UNKNOWN_FIELD = 'unknown field "surname" in proto.GreetRequest\n'

    const greetRequestWithInvalidField = {
      uponReceiving: 'a greeting',
      withRequest: {
        method: 'POST',
        path: '/api/v0/greet',
        body: REQUEST_DATA_INVALID_FIELD
      }
    }

    const preconditionFailedResponse = {
      status: 412,
      headers: {
        'Content-Type': 'text/plain; charset=utf-8'
      },
      body: RESPONSE_DATA_UNKNOWN_FIELD
    }

    // add expectations
    it('returns a successful body', async () => {
      const interaction = {
        state: 'i am greeted with my name',
        ...greetRequest,
        willRespondWith: successResponse
      }
      await provider.addInteraction(interaction)

      return Hello_Greet({
        $domain: provider.mockService.baseUrl,
        body: REQUEST_DATA
      }).then(response => {
        expect(response.data).toEqual(RESPONSE_DATA)
      })
    })

    it('returns an error with a wrong field', async () => {
      const interaction = {
        state: 'i have sent a wrong field name',
        ...greetRequestWithInvalidField,
        willRespondWith: preconditionFailedResponse
      }
      await provider.addInteraction(interaction)

      return Hello_Greet({
        $domain: provider.mockService.baseUrl,
        body: REQUEST_DATA_INVALID_FIELD
      }).catch((error) => {
        expect(error).toHaveProperty('response.status', 412)
        expect(error).toHaveProperty('response.data', RESPONSE_DATA_UNKNOWN_FIELD)
      })
    })
  })
})
