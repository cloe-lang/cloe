Feature: Conversion
  Scenario: Convert a dictionary to a list
    Given a file named "main.cloe" with:
    """
    (write (toList {123 456 "foo" "bar"}))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    [[123 456] ["foo" "bar"]]
    """

  Scenario: Convert a list to a list
    Given a file named "main.cloe" with:
    """
    (write (toList [123 nil 456 "foo" true "bar" false]))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    [123 nil 456 "foo" true "bar" false]
    """

  Scenario: Convert a string to a list
    Given a file named "main.cloe" with:
    """
    (write (toList "Cloe is good."))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    ["C" "l" "o" "e" " " "i" "s" " " "g" "o" "o" "d" "."]
    """
