Feature: Let statement
  Scenario: Define a variable
    Given a file named "main.tisp" with:
    """
    (let foo 123)
    (write foo)
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    123
    """

  Scenario: Define a function
    Given a file named "main.tisp" with:
    """
    (let (f x) x)
    (write (f 42))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    42
    """

  Scenario: Define a variable in a function
    Given a file named "main.tisp" with:
    """
    (let (foo x)
      (let bar (+ x x))
      bar)

    (write (foo 21))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    42
    """

  Scenario: Define nested functions
    Given a file named "main.tisp" with:
    """
    (let (f x)
      (let (g y) (+ x y))
      (g 42))

    (write (f 2007))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    2049
    """

  Scenario: Define complex nested functions
    Given a file named "main.tisp" with:
    """
    (let (foo x y)
      (let bar (+ x x))
      (let baz (- x y))
      (* bar baz (+ x y)))

    (write (foo 2 3))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    -20
    """

  Scenario: Define a function with a nested function definition
    Given a file named "main.tisp" with:
    """
    (let (f x)
      (let (g y)
        (let (h z)
          (+ x y z))
        h)
      ((g 456) 789))

    (write (f 123))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    1368
    """

  Scenario: Define a variable shadowing another
    Given a file named "main.tisp" with:
    """
    (let (f x)
      (let (g x) x)
      (g 42))

    (write (f 123456))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    42
    """
