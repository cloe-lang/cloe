Feature: Function application
  Scenario: Apply a function to a positional argument
    Given a file named "main.cloe" with:
    """
    (def (f x) x)
    (write (f 123))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    123
    """

  Scenario: Apply a function to 2 positional arguments
    Given a file named "main.cloe" with:
    """
    (def (f x y) (+ x y))
    (write (f 123 456))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    579
    """

  Scenario: Apply a function to complex arguments
    Given a file named "main.cloe" with:
    """
    (def (f x ..args . foo 4 ..kwargs) (+ x ..args foo))
    (write (f 1 2 . foo 3))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    6
    """
