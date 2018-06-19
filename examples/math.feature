Feature: Math
  Scenario: Add 2 numbers
    Given a file named "main.cloe" with:
    """
    (write (+ 2016 33))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "2049"

  Scenario: Subtract a number from the other
    Given a file named "main.cloe" with:
    """
    (write (- 2049 33))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "2016"

  Scenario: Divide a number by the other
    Given a file named "main.cloe" with:
    """
    (write (/ 84 2))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "42"

  Scenario: Use a negative number literal
    Given a file named "main.cloe" with:
    """
    (write (- 2007 -42))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "2049"
