Feature: I/O
  Scenario: Print a string
    Given a file named "main.cloe" with:
    """
    (print "Hello!")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello!"

  Scenario: Print a string with end argument
    Given a file named "main.cloe" with:
    """
    (print "Hello!" . end "!!")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello!!!"

  Scenario: Print multiple strings
    Given a file named "main.cloe" with:
    """
    (print "Hello," "world!")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Print a number
    Given a file named "main.cloe" with:
    """
    (print 42)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "42"

  Scenario: Print a nil
    Given a file named "main.cloe" with:
    """
    (print nil)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "nil"

  Scenario: Print multiple arguments of different types
    Given a file named "main.cloe" with:
    """
    (print "string" 42 nil true)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "string 42 nil true"

  Scenario: Print a string to stderr
    Given a file named "main.cloe" with:
    """
    (print "This is stderr." . file 2)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""
    And the stderr should contain exactly "This is stderr."

  Scenario: Print a string to a file
    Given a file named "main.cloe" with:
    """
    (print "This is content." . file "output.txt")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""
    And the file "output.txt" should contain exactly:
    """
    This is content.
    """

  Scenario: Print with a wrong file argument
    Given a file named "main.cloe" with:
    """
    (print 42 . file nil)
    """
    When I run `cloe main.cloe`
    Then the exit status should not be 0
    And the stdout should contain exactly ""
    And the stderr should contain "Error"
