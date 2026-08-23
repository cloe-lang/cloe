Feature: Recursion
  Scenario: Define a recursive function
    Given a file named "main.cloe" with:
    """
    (def (factorial n)
         (if (= n 0) 1 (* n (factorial (- n 1)))))

    (print (factorial 5))
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

    (print ((createFactorial) 5))
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
      (print (even? 0))
      (print (odd? 0))
      (print (even? 1))
      (print (odd? 1))
      (print (even? 2))
      (print (odd? 2))
      (print (even? 3))
      (print (odd? 3))
      (print (even? 42))
      (print (odd? 42))
      (print (even? 2049))
      (print (odd? 2049)))
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
      (print (even? 0))
      (print (odd? 0))
      (print (even? 1))
      (print (odd? 1))
      (print (even? 2))
      (print (odd? 2))
      (print (even? 3))
      (print (odd? 3))
      (print (even? 42))
      (print (odd? 42))
      (print (even? 2049))
      (print (odd? 2049)))
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
      (print (even? 0))
      (print (odd? 0))
      (print (even? 1))
      (print (odd? 1))
      (print (even? 2))
      (print (odd? 2))
      (print (even? 3))
      (print (odd? 3))
      (print (even? 42))
      (print (odd? 42))
      (print (even? 2049))
      (print (odd? 2049)))
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

    (let even? (@ (foo) 1))
    (let odd? (@ (foo) 2))

    (seq!
      (print (even? 0))
      (print (odd? 0))
      (print (even? 1))
      (print (odd? 1))
      (print (even? 2))
      (print (odd? 2))
      (print (even? 3))
      (print (odd? 3))
      (print (even? 42))
      (print (odd? 42))
      (print (even? 2049))
      (print (odd? 2049)))
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
