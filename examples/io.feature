Feature: IO
  Scenario: Read stdin
    Given a file named "main.cloe" with:
    """
    (write (read) . end "")
    """
    And a file named "test.txt" with:
    """
    foo
    bar
    baz
    """
    When I run the following commands:
    """
    cloe main.cloe < test.txt
    """
    Then the stdout should contain exactly:
    """
    foo
    bar
    baz
    """

  Scenario: Read a file
    Given a file named "main.cloe" with:
    """
    (write (read . file "test.txt") . end "")
    """
    And a file named "test.txt" with:
    """
    foo
    bar
    baz
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly:
    """
    foo
    bar
    baz
    """

  Scenario: Write a string
    Given a file named "main.cloe" with:
    """
    (write "Hello!")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello!"

  Scenario: Write a string with end argument
    Given a file named "main.cloe" with:
    """
    (write "Hello!" . end "!!")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello!!!"

  Scenario: Write multiple strings
    Given a file named "main.cloe" with:
    """
    (write "Hello," "world!")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "Hello, world!"

  Scenario: Write a number
    Given a file named "main.cloe" with:
    """
    (write 42)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "42"

  Scenario: Write a nil
    Given a file named "main.cloe" with:
    """
    (write nil)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "nil"

  Scenario: Write multiple arguments of different types
    Given a file named "main.cloe" with:
    """
    (write "string" 42 nil true)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly "string 42 nil true"

  Scenario: Write a string to stderr
    Given a file named "main.cloe" with:
    """
    (write "This is stderr." . file 2)
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""
    And the stderr should contain exactly "This is stderr."

  Scenario: Write a string to a file
    Given a file named "main.cloe" with:
    """
    (write "This is content." . file "output.txt")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""
    And the file "output.txt" should contain exactly:
    """
    This is content.
    """

  Scenario: Write with a wrong file argument
    Given a file named "main.cloe" with:
    """
    (write 42 . file nil)
    """
    When I run `cloe main.cloe`
    Then the exit status should not be 0
    And the stdout should contain exactly ""
    And the stderr should contain "Error"
