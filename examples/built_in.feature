Feature: Built-in functions
  Scenario: Add 2 numbers
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
