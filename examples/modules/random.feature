Feature: Random number generator
  Scenario: Import random module
    Given a file named "main.cloe" with:
    """
    (import "random")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""

  Scenario: Generate random numbers
    Given a file named "main.cloe" with:
    """
    (import "random")

    (seq! (print (<= 0 (random.number) 1)))
    """
    When I run `cloe main.cloe`
    Then the stdout should contain exactly "true"
