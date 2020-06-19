Feature: Hello

	Scenario: Greet
		Given user "user1" has been created with default attributes
		And user "user1" has logged in using the webUI
		And the user browses to the hello page
		When the user submits "User One" to the Greet input
		Then "Hello User One" should be shown in the hello screen
