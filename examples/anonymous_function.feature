Feature: Anonymous function
  Scenario: Call an anonymous function
    Given a file named "main.tisp" with:
    """
    (write ((\ (x) x) "Hello, world!"))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly "Hello, world!"
