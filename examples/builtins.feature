Feature: Built-in functions
  Scenario: Get types of values
    Given a file named "main.cloe" with:
    """
    (seq!
      (write (typeOf true))
      (write (typeOf {"key" "value"}))
      (write (typeOf []))
      (write (typeOf nil))
      (write (typeOf 42))
      (write (typeOf "foo"))
      (write (typeOf +))
      (write (typeOf (partial + 1))))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    boolean
    dictionary
    list
    nil
    number
    string
    function
    function
    """

  Scenario: Map a function to a list
    Given a file named "main.cloe" with:
    """
    (write (map (\ (x) (* x x)) [1 2 3]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    [1 4 9]
    """

  Scenario: Calculate indices of elements in a list
    Given a file named "main.cloe" with:
    """
    (let l [1 2 3 42 -3 "foo"])
    (seq!
      (write (index l 42))
      (write (index l 2))
      (write (index l "foo")))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    4
    2
    6
    """

  Scenario: Use multiple conditions with if function
    Given a file named "main.cloe" with:
    """
    (def (no) (write "No"))

    (if false no true (write "Yes") false no no)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    Yes
    """

  Scenario: Use boolean operators
    Given a file named "main.cloe" with:
    """
    (seq!
      (write (not true))
      (write (not false))
      (write (and true))
      (write (or true))
      (write (and true false))
      (write (or true false))
      (write (and true true))
      (write (or false false))
      (write (and true false true))
      (write (or true false true)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    false
    true
    true
    true
    false
    true
    true
    false
    false
    true
    """

  Scenario: Slice lists
    Given a file named "main.cloe" with:
    """
    (seq!
      (write (slice [1 2 3]))
      (write (slice [1 2 3] . start 1))
      (write (slice [1 2 3] . start 2))
      (write (slice [1 2 3] . end 1))
      (write (slice [1 2 3] . start 3))
      (write (slice [1 2 3] . start 4))
      (write (slice [1 2 3] . start 5))
      (write (slice [1 2 3] . start 2 end 3))
      (write (slice [1 2 3] . start 1 end 2)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    [1 2 3]
    [1 2 3]
    [2 3]
    [1]
    [3]
    []
    []
    [2 3]
    [1 2]
    """

  Scenario: Slice an infinite list
    Given a file named "main.cloe" with:
    """
    (def (f) [42 ..(f)])

    (write (slice (f) . start 1 end 3))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    [42 42 42]
    """

  Scenario: Slice strings
    Given a file named "main.cloe" with:
    """
    (seq! ..(map (\ (x) (write (dump x))) [
        (slice "abc")
        (slice "abc" . start 1)
        (slice "abc" . start 2)
        (slice "abc" . end 1)
        (slice "abc" . start 3)
        (slice "abc" . start 4)
        (slice "abc" . start 5)
        (slice "abc" . start 2 end 3)
        (slice "abc" . start 1 end 2)
      ]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    "abc"
    "abc"
    "bc"
    "a"
    "c"
    ""
    ""
    "bc"
    "ab"
    """

  Scenario: Calculate maximum and minimum of numbers
    Given a file named "main.cloe" with:
    """
    (seq!
      (write (max 1))
      (write (max 1 2))
      (write (max 1 2 3))
      (write (max 3))
      (write (max 3 2))
      (write (max 3 2 1))
      (write (max 3 2 4 -3 123 -45 1))
      (write (min 1))
      (write (min 1 2))
      (write (min 1 2 3))
      (write (min 3))
      (write (min 3 2))
      (write (min 3 2 1))
      (write (min 3 2 4 -3 123 -45 1)))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    1
    2
    3
    3
    3
    3
    123
    1
    1
    1
    3
    2
    1
    -45
    """

  Scenario: Zip lists
    Given a file named "main.cloe" with:
    """
    (write (zip [1 2 3] ["foo" "bar" "baz"]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    [[1 "foo"] [2 "bar"] [3 "baz"]]
    """

  Scenario: Check if values are ordered or not
    Given a file named "main.cloe" with:
    """
    (seq!
      ..(map (\ (x) (write (ordered? x))) [
        123
        "foo"
        []
        [123]
        nil
        true
        {}
        [{}]
        [123 {}]
      ]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    true
    true
    true
    true
    false
    false
    false
    false
    false
    """

  Scenario: Check if values are ordered or not
    Given a file named "main.cloe" with:
    """
    (seq!
      ..(map write [
        (boolean? true)
        (boolean? 42)
        (dictionary? {"foo" 42})
        (dictionary? "foo")
        (function? (\ (x) x))
        (function? [])
        (list? [42 "foo"])
        (list? nil)
        (nil? nil)
        (nil? "foo")
        (number? 42)
        (number? [])
        (string? "foo")
        (string? nil)
      ]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    true
    false
    true
    false
    true
    false
    true
    false
    true
    false
    true
    false
    true
    false
    """
