Feature: Match expression
  Scenario: Match scalars
    Given a file named "main.tisp" with:
    """
    (write (match 42 42 "Matched!"))
    (write (match 42
      2049 "Not matched..."
      42 "Matched!"))

    (write (match "Tisp" "Tisp" "Matched!"))
    (write (match "Tisp"
      "tisp" "Not matched..."
      "Tisp" "Matched!"))

    (write (match true true "Matched!"))
    (write (match true
      false "Not matched..."
      true "Matched!"))

    (write (match nil nil "Matched!"))

    (write (match "Matched!" x x))
    """
    When I successfully run `tisp main.tisp`
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
    Given a file named "main.tisp" with:
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
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    Matched!
    Matched!
    Matched!
    Matched!
    """

  Scenario: Use wildcards
    Given a file named "main.tisp" with:
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
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    Matched!
    Matched!
    Matched!
    """

  Scenario: Nest collections
    Given a file named "main.tisp" with:
    """
    (write (match {"foo" 42 "bar" ["The pattern" "is" "Matched!"]}
      {"bar" [foo "is" baz] "foo" 42} baz
      {"foo" foo "bar" bar} "Not matched..."))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    Matched!
    """
