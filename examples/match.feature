Feature: Match expression
  Scenario: Match scalars
    Given a file named "main.coel" with:
    """
    (write (match 42 42 "Matched!"))
    (write (match 42
      2049 "Not matched..."
      42 "Matched!"))

    (write (match "Coel" "Coel" "Matched!"))
    (write (match "Coel"
      "coel" "Not matched..."
      "Coel" "Matched!"))

    (write (match true true "Matched!"))
    (write (match true
      false "Not matched..."
      true "Matched!"))

    (write (match nil nil "Matched!"))

    (write (match "Matched!" x x))
    """
    When I successfully run `coel main.coel`
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
    Given a file named "main.coel" with:
    """
    (write (match [] [] "Matched!"))
    (write (match [42]
      [] "Not matched..."
      [42 42] "Not matched..."
      [42] "Matched!"))

    (write (match {} {} "Matched!"))
    (write (match {"foo" 42}
      {} "Not matched..."
      {"foo" 2049} "Not matched..."
      {"foo" 42 "bar" 2049} "Not matched..."
      {"bar" 42} "Not matched..."
      {"foo" 42} "Matched!"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Matched!
    Matched!
    Matched!
    Matched!
    """

  Scenario: Use wildcards
    Given a file named "main.coel" with:
    """
    (write (match "Matched!"
      42 "Not matched..."
      x x))

    (write (match [42 2049]
      [] "Not matched..."
      [2049] "Not matched..."
      [42 42] "Not matched..."
      [foo 42] "Not matched..."
      [42 bar 2049] "Not matched..."
      [foo 2049] "Matched!"))

    (write (match {"foo" 42 "bar" "Matched!"}
      {} "Not matched..."
      {"foo" 42} "Not matched..."
      {"bar" 42} "Not matched..."
      {"foo" foo "bar" 42} "Not matched..."
      {"bar" bar "foo" 2049} "Not matched..."
      {"foo" foo "bar" bar} bar))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Matched!
    Matched!
    Matched!
    """

  Scenario: Nest collections
    Given a file named "main.coel" with:
    """
    (write (match {"foo" 42 "bar" ["The pattern" "is" "Matched!"]}
      {"bar" [foo "is" baz] "foo" 42} baz
      {"foo" foo "bar" bar} "Not matched..."))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Matched!
    """

  Scenario: Use rest pattern of list
    Given a file named "main.coel" with:
    """
    (write (match [[42 2049] ["This" "is" "Matched!"]]
      [[..foo] ["This" "is" "not" "Matched!"]] "Not matched..."
      [[..foo] ["This" ..bar]] (bar 1)))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Matched!
    """

  Scenario: Use rest pattern of dictionary
    Given a file named "main.coel" with:
    """
    (write (match {"foo" {42 2049} "bar" {"This" "Matched!"}}
      [[..foo] ["This" "is" "not" "Matched!"]] "Not matched..."
      {"foo" {42 2050} "bar" {"This" "Matched!"}} "Not matched..."
      {"foo" {..foo} "bar" {"this" bar}} "Not matched..."
      {"foo" {..foo} "bar" {"This" bar}} bar))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Matched!
    """

  Scenario: Match an invalid list with a valid pattern
    Given a file named "main.coel" with:
    """
    (write (match (prepend "Matched!" {})
      [x ..xs] x))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    Matched!
    """

  Scenario: Use match expression with let statements
    Given a file named "main.coel" with:
    """
    (let y (match [123 456 789]
      [x ..xs] xs))

    (write y)
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    [456 789]
    """

  Scenario: Match collections with patterns in let statements
    Given a file named "main.coel" with:
    """
    (let [x y ..xs] ["foo" "bar" "baz"])
    (let {"foo" value ..rest} {"foo" 42 "bar" 2049})

    (seq
      (write y)
      (write value))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    bar
    42
    """

  Scenario: Use let-match statements in function definitions
    Given a file named "main.coel" with:
    """
    (def (f x y)
      (let [_ x ..xs] x)
      (let {"foo" y ..rest} y)
      [x y])

    (write (f ["foo" "bar" "baz"] {"foo" 42 "bar" 2049}))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    ["bar" 42]
    """
