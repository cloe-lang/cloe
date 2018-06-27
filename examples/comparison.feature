Feature: Comparison
  Scenario: Use < operator
    Given a file named "main.cloe" with:
    """
    (print (if (< 1 2 3) "OK" "Not OK"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "OK"

  Scenario: Use <= operator
    Given a file named "main.cloe" with:
    """
    (print (if (<= 1 1 3) "OK" "Not OK"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "OK"

  Scenario: Use > operator
    Given a file named "main.cloe" with:
    """
    (print (if (> 3 2 1) "OK" "Not OK"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "OK"

  Scenario: Use >= operator
    Given a file named "main.cloe" with:
    """
    (print (if (>= 3 1 1) "OK" "Not OK"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "OK"

  Scenario: Cannot use < operator for boolean values
    Given a file named "main.cloe" with:
    """
    (print (< false true))
    """
    When I run `cloe main.cloe`
    Then the exit status should not be 0
