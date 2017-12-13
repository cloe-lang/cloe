Feature: Comparison
  Scenario: Use < operator
    Given a file named "main.coel" with:
    """
    (write (if (< 1 2 3) "OK" "Not OK"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    OK
    """

  Scenario: Use <= operator
    Given a file named "main.coel" with:
    """
    (write (if (<= 1 1 3) "OK" "Not OK"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    OK
    """

  Scenario: Use > operator
    Given a file named "main.coel" with:
    """
    (write (if (> 3 2 1) "OK" "Not OK"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    OK
    """

  Scenario: Use >= operator
    Given a file named "main.coel" with:
    """
    (write (if (>= 3 1 1) "OK" "Not OK"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    OK
    """

  Scenario: Cannot use < operator for boolean values
    Given a file named "main.coel" with:
    """
    (write (< false true))
    """
    When I run `coel main.coel`
    Then the exit status should not be 0
