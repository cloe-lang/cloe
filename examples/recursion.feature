Feature: Recursion
  Scenario: Define a recursive function
    Given a file named "main.cloe" with:
    """
    (def (factorial n)
         (if (= n 0) 1 (* n (factorial (- n 1)))))

    (write (factorial 5))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "120"

  Scenario: Define a recursive function in a function definition
    Given a file named "main.cloe" with:
    """
    (def (createFactorial)
         (def (factorial n)
              (if (= n 0) 1 (* n (factorial (- n 1)))))
         factorial)

    (write ((createFactorial) 5))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "120"

  Scenario: Define 2 mutually recursive functions
    Given a file named "main.cloe" with:
    """
    (mr
      (def (even? n)
        (if (= n 0) true (odd? (- n 1))))
      (def (odd? n)
        (if (= n 0) false (even? (- n 1)))))

    (seq!
      (write (even? 0))
      (write (odd? 0))
      (write (even? 1))
      (write (odd? 1))
      (write (even? 2))
      (write (odd? 2))
      (write (even? 3))
      (write (odd? 3))
      (write (even? 42))
      (write (odd? 42))
      (write (even? 2049))
      (write (odd? 2049)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    true
    false
    false
    true
    true
    false
    false
    true
    true
    false
    false
    true
    """

  Scenario: Define mutually recursive functions with a nested let statement
    Given a file named "main.cloe" with:
    """
    (mr
      (def (even? n)
        (let o? odd?)
        (if (= n 0) true (o? (- n 1))))
      (def (odd? n)
        (if (= n 0) false (even? (- n 1)))))

    (seq!
      (write (even? 0))
      (write (odd? 0))
      (write (even? 1))
      (write (odd? 1))
      (write (even? 2))
      (write (odd? 2))
      (write (even? 3))
      (write (odd? 3))
      (write (even? 42))
      (write (odd? 42))
      (write (even? 2049))
      (write (odd? 2049)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    true
    false
    false
    true
    true
    false
    false
    true
    true
    false
    false
    true
    """

  Scenario: Define mutually recursive functions with a shadowed variable
    Given a file named "main.cloe" with:
    """
    (mr
      (def (even? n)
        (let even? odd?)
        (if (= n 0) true (even? (- n 1))))
      (def (odd? n)
        (if (= n 0) false (even? (- n 1)))))

    (seq!
      (write (even? 0))
      (write (odd? 0))
      (write (even? 1))
      (write (odd? 1))
      (write (even? 2))
      (write (odd? 2))
      (write (even? 3))
      (write (odd? 3))
      (write (even? 42))
      (write (odd? 42))
      (write (even? 2049))
      (write (odd? 2049)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    true
    false
    false
    true
    true
    false
    false
    true
    true
    false
    false
    true
    """

  Scenario: Define mutually recursive functions in a function
    Given a file named "main.cloe" with:
    """
    (def (foo)
      (mr
        (def (even? n) (if (= n 0) true  (odd?  (- n 1))))
        (def (odd? n)  (if (= n 0) false (even? (- n 1)))))
      [even? odd?])

    (let even? ((foo) 1))
    (let odd? ((foo) 2))

    (seq!
      (write (even? 0))
      (write (odd? 0))
      (write (even? 1))
      (write (odd? 1))
      (write (even? 2))
      (write (odd? 2))
      (write (even? 3))
      (write (odd? 3))
      (write (even? 42))
      (write (odd? 42))
      (write (even? 2049))
      (write (odd? 2049)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    true
    false
    false
    true
    true
    false
    false
    true
    true
    false
    false
    true
    """
