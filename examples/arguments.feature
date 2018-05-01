Feature: Arguments
  Scenario: Define a recursive function
    Given a file named "main.cloe" with:
    """
    (def (func x1 x2 ..args . y1 0 y2 1 ..kwargs)
         (+ x1 x2 ..args y1 y2))
    (let foo 1)

    (write (func 1 1 1 ..[1 foo foo] . y1 1 foo 100000000 ..{"y2" 1}))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    8
    """

  Scenario: Override keyword arguments
    Given a file named "main.cloe" with:
    """
    (def (func . x nil) x)

    (seq!
      (write (func . x nil ..{"x" 42}))
      (write (func . ..{"x" nil} x 42))
      (write (func . ..{"x" 42} ..{"x" nil} x 42)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    42
    42
    42
    """
