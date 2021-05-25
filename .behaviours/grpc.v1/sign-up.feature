Feature: Sign Up
  Background: Setup of System Under Test (SUT)
    Given a running "Identity" service

  Scenario: Happy case
    When I create a "sign-up" request
    And set "username" to "alice"
    And set "password" to "0123456789"
    And send the "sign-up" request to the "Identity" service
    Then the "error" should be ""

  Scenario: Username already exists
    Given a "signed-up" user: "alice:+123456789"
    When I create a "sign-up" request
    And set "username" to "alice"
    And set "password" to "_123456789"
    And send the "sign-up" request to the "Identity" service
    Then the "error" should be ""

  Scenario Outline: Bad input
    When I create a "sign-up" request
    And set "username" to "<username>"
    And set "password" to "<password>"
    And send the "sign-up" request to the "Identity" service
    Then the "error.code" should be <code>
    And the "error.message" should be <message>
    Examples:
      | username | password   | code              | message              | comment                       |
      |          | 0123456789 | "InvalidArgument" | "empty_username"     | Empty username is not allowed |
      | alice    | 012345678  | "InvalidArgument" | "too_short_password" | Password requirement is >= 10 |
