Feature: Conversion
  Scenario: Convert a dictionary to a list
    Given a file named "main.coel" with:
    """
    (write (toList {123 456 "foo" "bar"}))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    [[123 456] ["foo" "bar"]]
    """

  Scenario: Convert a list to a list
    Given a file named "main.coel" with:
    """
    (write (toList [123 nil 456 "foo" true "bar" false]))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    [123 nil 456 "foo" true "bar" false]
    """

  Scenario: Convert a string to a list
    Given a file named "main.coel" with:
    """
    (write (toList "Coel is good."))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    ["C" "o" "e" "l" " " "i" "s" " " "g" "o" "o" "d" "."]
    """
