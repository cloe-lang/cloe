Feature: Anonymous function
  Scenario: Call an anonymous function
    Given a file named "main.coel" with:
    """
    (write ((\ (x) x) "Hello, world!"))
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly "Hello, world!"
