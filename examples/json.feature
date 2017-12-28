Feature: JSON
  Scenario: Import JSON module
    Given a file named "main.coel" with:
    """
    (import "json")
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly ""

  Scenario: Decode strings in JSON
    Given a file named "main.coel" with:
    """
    (import "json")

    (write (json.decode "{\"foo\": 42}"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    {"foo" 42}
    """

  Scenario: Encode values into JSON
    Given a file named "main.coel" with:
    """
    (import "json")

    (seq!
      (write (json.encode {"foo" 42}))
      (write (json.encode {123 nil}))
      (write (json.encode {nil "bar"})))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    {"foo":42}
    {"123":null}
    {"null":"bar"}
    """
