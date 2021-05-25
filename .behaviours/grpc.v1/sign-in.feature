Feature: Sign In
  Background: Setup of System Under Test (SUT)
    Given a running "Identity" service
    And a "signed-up" user: "alice:0123456789"

  Scenario: Happy case
    When I create a "sign-in" request
    And set "username" to "alice"
    And set "password" to "0123456789"
    And send the "sign-in" request to the "Identity" service
    Then the "error" should be ""
    And the "access_token" should not be ""

  Scenario Outline: Bad input
    When I create a "sign-in" request
    And set "username" to "<username>"
    And set "password" to "<password>"
    And send the "sign-in" request to the "Identity" service
    Then the "error.code" should be <code>
    And the "error.message" should be <message>
    Examples:
      | username | password   | code              | message                      | comment |
      |          | 0123456789 | "InvalidArgument" | "wrong_username_or_password" |         |
      | alice    | 012345678  | "InvalidArgument" | "wrong_username_or_password" |         |
