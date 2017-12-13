Feature: Built-in functions
  Scenario: Get types of values
    Given a file named "main.coel" with:
    """
    (seq
      (write (typeOf true))
      (write (typeOf {"key" "value"}))
      (write (typeOf []))
      (write (typeOf nil))
      (write (typeOf 42))
      (write (typeOf "foo"))
      (write (typeOf +))
      (write (typeOf (partial + 1))))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    bool
    dict
    list
    nil
    number
    string
    function
    function
    """

  Scenario: Map a function to a list
    Given a file named "main.coel" with:
    """
    (write (map (\ (x) (* x x)) [1 2 3]))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    [1 4 9]
    """

  Scenario: Calculate indices of elements in a list
    Given a file named "main.coel" with:
    """
    (let l [1 2 3 42 -3 "foo"])
    (seq
      (write (indexOf l 42))
      (write (indexOf l 2))
      (write (indexOf l "foo")))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    3
    1
    5
    """
