Feature: Operating System
  Scenario: Import OS module
    Given a file named "main.coel" with:
    """
    (import "os")
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly ""

  Scenario: Exit with 0
    Given a file named "main.coel" with:
    """
    (import "os")

    (os.exit)
    """
    When I run `coel main.coel`
    Then the exit status should be 0

  Scenario: Exit with 1
    Given a file named "main.coel" with:
    """
    (import "os")

    (os.exit 1)
    """
    When I run `coel main.coel`
    Then the exit status should be 1
