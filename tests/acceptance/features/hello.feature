Feature: hello API

  As I user
  I want to be greeted by the greeting API with my name
  So that I can feel important

  Scenario Outline: response on correct POST messages
    When sending 'POST' to 'api/v0/greet' with
      | body | {"name": "<name>"} |
    Then response http status code should be 201
    And response body should be
      """
      {"message":"Hello <name>"}

      """
    Examples:
      | name         | comment                |
      | Peter        | simple ASCII character |
      | आर्तुर       | UTF caracters          |
      | 🐅🌵🍉🥐🐊🐅 | emojis                 |

  Scenario Outline: response should be also correct when various headers are set
    When sending 'POST' to 'api/v0/greet' with
      | body    | {"name": "Peter"} |
      | headers | <headers>         |
    Then response http status code should be 201
    And response body should be
      """
      {"message":"Hello Peter"}

      """
    Examples:
      | headers                                     |
      | {"XMyHEader": "what a rubbish"}             |
      | {"Authorization": "Basic YXJ0dXI6YXJ0dXIK"} |


  Scenario Outline: response on GET requests on different URLS
    When sending 'GET' to 'api/v0/greet'
    Then response http status code should be 405
    Examples:
      | url          |
      | api/v0/greet |
      | api/v1/greet |
      | api/v0/hello |
      | greet        |
      | api          |
      | api/v0       |
