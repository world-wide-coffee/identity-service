Feature: Change Password
  Background: Setup of System Under Test (SUT)
    Given a running "Identity" service
    And a "signed-in" user: "alice:0123456789"

  Scenario: Happy case
    When I create a "change-password" request
    And set "access_token" to "{{sign-in.reply.access_token}}"
    And set "id" to "{{sign-up.reply.id}}"
    And set "old_password" to "0123456789"
    And set "new_password" to "9876543210"
    And send the "change-password" request to the "Identity" service
    Then the "error" should be ""

    When I "sign-in" with "alice:9876543210"
    Then the "error" should be ""
    And the "access_token" should not be ""