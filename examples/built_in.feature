Feature: Built-in functions
  Scenario: Get types of values
    Given a file named "main.tisp" with:
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
    When I successfully run `tisp main.tisp`
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
    Given a file named "main.tisp" with:
    """
    (write (map (\ (x) (* x x)) [1 2 3]))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    [1 4 9]
    """
