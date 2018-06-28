Feature: Functions
  Scenario: Define a function
    Given a file named "main.cloe" with:
    """
    (def (f x) x)
    (print (f 42))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "42"

  Scenario: Apply a function to a positional argument
    Given a file named "main.cloe" with:
    """
    (def (f x) x)
    (print (f 123))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "123"

  Scenario: Apply a function to 2 positional arguments
    Given a file named "main.cloe" with:
    """
    (def (f x y) (+ x y))
    (print (f 123 456))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "579"

  Scenario: Override keyword arguments
    Given a file named "main.cloe" with:
    """
    (def (func . x nil) x)

    (seq!
      (print (func . x nil ..{"x" 42}))
      (print (func . ..{"x" nil} x 42))
      (print (func . ..{"x" 42} ..{"x" nil} x 42)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    42
    42
    42
    """

  Scenario: Apply a function to complex arguments
    Given a file named "main.cloe" with:
    """
    (def (f x ..args . foo 4 ..kwargs) (+ x ..args foo))
    (print (f 1 2 . foo 3))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "6"

  Scenario: Apply a function to very complex arguments
    Given a file named "main.cloe" with:
    """
    (def (func x1 x2 ..args . y1 0 y2 1 ..kwargs)
         (+ x1 x2 ..args y1 y2))

    (print (func 1 1 1 ..[1 1 1] . y1 1 1 100000000 ..{"y2" 1}))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "8"

  Scenario: Define a variable
    Given a file named "main.cloe" with:
    """
    (let foo 123)
    (print foo)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "123"

  Scenario: Define a recursive variable
    Given a file named "main.cloe" with:
    """
    (let l [42 ..l])

    (print (l 1))
    (print (l 2))
    (print (l 3))
    """
    When I run `cloe main.cloe`
    Then the exit status should not be 0

  Scenario: Define a variable in a function
    Given a file named "main.cloe" with:
    """
    (def (foo x)
      (let bar (+ x x))
      bar)

    (print (foo 21))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "42"

  Scenario: Define multiple variables in a function
    Given a file named "main.cloe" with:
    """
    (def (foo x y)
      (let bar (+ x x))
      (let baz (- x y))
      (* bar baz (+ x y)))

    (print (foo 2 3))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "-20"

  Scenario: Define nested functions
    Given a file named "main.cloe" with:
    """
    (def (f x)
      (def (g y) (+ x y))
      (g 42))

    (print (f 2007))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "2049"

  Scenario: Define a deeply nested function
    Given a file named "main.cloe" with:
    """
    (def (f x)
      (def (g y)
        (def (h z)
          (+ x y z))
        h)
      ((g 456) 789))

    (print (f 123))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "1368"

  Scenario: Shadow an argument
    Given a file named "main.cloe" with:
    """
    (def (f x)
      (let x (+ x 1))
      (let x (+ x 1))
      x)

    (print (f 1))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "3"

  Scenario: Define an argument shadowing another
    Given a file named "main.cloe" with:
    """
    (def (f x)
      (def (g x) x)
      (g 42))

    (print (f 123456))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "42"

  Scenario: Call an anonymous function
    Given a file named "main.cloe" with:
    """
    (print ((\ (x) x) "Hello, world!"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Use let expressions
    Given a file named "main.cloe" with:
    """
    (def (foo ..xs)
      (match xs
        [] nil
        [x ..xs] (let
          y (+ ..xs)
          z (- y 0)
          (let [v w] xs
            (+ v w x y z)))))

    (print (foo 2 3 4))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "23"
