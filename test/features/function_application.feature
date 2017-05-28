Feature: Function application
  Scenario: Apply a function to a positional argument
    Given a file named "main.tisp" with:
    """
    (let (f x) x)
    (write (f 123))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    123
    """

  Scenario: Apply a function to 2 positional arguments
    Given a file named "main.tisp" with:
    """
    (let (f x y) (+ x y))
    (write (f 123 456))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    579
    """

  Scenario: Apply a function to complex arguments
    Given a file named "main.tisp" with:
    """
    (let (f x (y 2) ..args . foo (bar 4) ..kwargs) (+ x y foo bar))
    (write (f 1 . foo 3))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    10
    """
