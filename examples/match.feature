Feature: Match expression
  Scenario: Match scalars
    Given a file named "main.cloe" with:
    """
    (print (match 42 42 "Matched!"))
    (print (match 42
      2049 "Not matched..."
      42 "Matched!"))

    (print (match "Cloe" "Cloe" "Matched!"))
    (print (match "Cloe"
      "cloe" "Not matched..."
      "Cloe" "Matched!"))

    (print (match true true "Matched!"))
    (print (match true
      false "Not matched..."
      true "Matched!"))

    (print (match nil nil "Matched!"))

    (print (match "Matched!" x x))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    Matched!
    Matched!
    Matched!
    Matched!
    Matched!
    Matched!
    Matched!
    Matched!
    """

  Scenario: Match collections
    Given a file named "main.cloe" with:
    """
    (print (match [] [] "Matched!"))
    (print (match [42]
      [] "Not matched..."
      [42 42] "Not matched..."
      [42] "Matched!"))

    (print (match {} {} "Matched!"))
    (print (match {"foo" 42}
      {} "Not matched..."
      {"foo" 2049} "Not matched..."
      {"foo" 42 "bar" 2049} "Not matched..."
      {"bar" 42} "Not matched..."
      {"foo" 42} "Matched!"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    Matched!
    Matched!
    Matched!
    Matched!
    """

  Scenario: Use wildcards
    Given a file named "main.cloe" with:
    """
    (print (match "Matched!"
      42 "Not matched..."
      x x))

    (print (match [42 2049]
      [] "Not matched..."
      [2049] "Not matched..."
      [42 42] "Not matched..."
      [foo 42] "Not matched..."
      [42 bar 2049] "Not matched..."
      [foo 2049] "Matched!"))

    (print (match {"foo" 42 "bar" "Matched!"}
      {} "Not matched..."
      {"foo" 42} "Not matched..."
      {"bar" 42} "Not matched..."
      {"foo" foo "bar" 42} "Not matched..."
      {"bar" bar "foo" 2049} "Not matched..."
      {"foo" foo "bar" bar} bar))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    Matched!
    Matched!
    Matched!
    """

  Scenario: Nest collections
    Given a file named "main.cloe" with:
    """
    (print (match {"foo" 42 "bar" ["The pattern" "is" "Matched!"]}
      {"bar" [foo "is" baz] "foo" 42} baz
      {"foo" foo "bar" bar} "Not matched..."))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Matched!"

  Scenario: Use rest pattern of list
    Given a file named "main.cloe" with:
    """
    (print (match [[42 2049] ["This" "is" "Matched!"]]
      [[..foo] ["This" "is" "not" "Matched!"]] "Not matched..."
      [[..foo] ["This" ..bar]] (@ bar 2)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Matched!"

  Scenario: Use rest pattern of dictionary
    Given a file named "main.cloe" with:
    """
    (print (match {"foo" {42 2049} "bar" {"This" "Matched!"}}
      [[..foo] ["This" "is" "not" "Matched!"]] "Not matched..."
      {"foo" {42 2050} "bar" {"This" "Matched!"}} "Not matched..."
      {"foo" {..foo} "bar" {"this" bar}} "Not matched..."
      {"foo" {..foo} "bar" {"This" bar}} bar))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Matched!"

  Scenario: Match an invalid list with a valid pattern
    Given a file named "main.cloe" with:
    """
    (print (match ["Matched!" ..{}]
      [x ..xs] x))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Matched!"

  Scenario: Use match expression with let statements
    Given a file named "main.cloe" with:
    """
    (let y (match [123 456 789]
      [x ..xs] xs))

    (print y)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "[456 789]"

  Scenario: Match collections with patterns in let statements
    Given a file named "main.cloe" with:
    """
    (let [x y ..xs] ["foo" "bar" "baz"])
    (let {"foo" value ..rest} {"foo" 42 "bar" 2049})

    (seq!
      (print y)
      (print value))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    bar
    42
    """

  Scenario: Use let-match statements in function definitions
    Given a file named "main.cloe" with:
    """
    (def (f x y)
      (let [_ x ..xs] x)
      (let {"foo" y ..rest} y)
      [x y])

    (print (f ["foo" "bar" "baz"] {"foo" 42 "bar" 2049}))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    ["bar" 42]
    """

  Scenario: Nest match expressions
    Given a file named "main.cloe" with:
    """
    (print (match [1 2 3]
      [x ..xs] (match xs
        [y ..ys] (+ x y))))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "3"

  Scenario: Use many match expressions in a function
    Given a file named "main.cloe" with:
    """
    (def (f ..xs)
      (let xs [1 2 3 ..xs])
      (match xs
        [] (if true (match xs [x] x [..xs] "OOOOK"))
        [x y ..xs] (match xs [] [x y] [x y ..xs] [..xs (+ x y)])))

    (print (f 4 5 6))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "[5 6 7]"

  Scenario: Use similar match expressions in a expression
    Given a file named "main.cloe" with:
    """
    (print [(match [1] [x] x) (match [1] [x] x)])
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "[1 1]"
