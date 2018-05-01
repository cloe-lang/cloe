Feature: Operating System
  Scenario: Import OS module
    Given a file named "main.cloe" with:
    """
    (import "os")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""

  Scenario: Exit with 0
    Given a file named "main.cloe" with:
    """
    (import "os")

    (os.exit)
    """
    When I run `cloe main.cloe`
    Then the exit status should be 0

  Scenario: Exit with 1
    Given a file named "main.cloe" with:
    """
    (import "os")

    (os.exit 1)
    """
    When I run `cloe main.cloe`
    Then the exit status should be 1
