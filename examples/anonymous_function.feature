Feature: Anonymous function
  Scenario: Call an anonymous function
    Given a file named "main.cloe" with:
    """
    (write ((\ (x) x) "Hello, world!"))
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"
