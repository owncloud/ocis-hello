Feature: Hello

	Scenario: Greet
		Given user "admin" has logged in using the webUI
		And the user browses to the hello page
		When the user submits "User One" to the Greet input
		Then "Hello User One" should be shown in the hello screen
