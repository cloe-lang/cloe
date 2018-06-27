Feature: JSON
  Scenario: Import JSON module
    Given a file named "main.cloe" with:
    """
    (import "json")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""

  Scenario: Decode strings in JSON
    Given a file named "main.cloe" with:
    """
    (import "json")

    (print (json.decode "{\"foo\": 42}"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    {"foo" 42}
    """

  Scenario: Encode values into JSON
    Given a file named "main.cloe" with:
    """
    (import "json")

    (seq!
      (print (json.encode {"foo" 42}))
      (print (json.encode {123 nil}))
      (print (json.encode {nil "bar"})))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    {"foo":42}
    {"123":null}
    {"null":"bar"}
    """
