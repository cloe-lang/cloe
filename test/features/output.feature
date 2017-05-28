Feature: Output
  Scenario: Evaluate multiple outputs
    Given a file named "main.tisp" with:
    """
    (write 123)
    (write 456)
    (write 789)
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain:
    """
    123
    """
    And the stdout should contain:
    """
    456
    """
    And the stdout should contain:
    """
    789
    """

  Scenario: Evaluate an expanded output
    Given a file named "main.tisp" with:
    """
    ..[(write 123) (write 456) (write 789)]
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain:
    """
    123
    """
    And the stdout should contain:
    """
    456
    """
    And the stdout should contain:
    """
    789
    """
