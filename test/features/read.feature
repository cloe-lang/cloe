Feature: Read function
  Scenario: Read stdin
    Given a file named "main.tisp" with:
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
    tisp main.tisp < test.txt
    """
    Then the stdout should contain exactly:
    """
    foo
    bar
    baz
    """

  Scenario: Read a file
    Given a file named "main.tisp" with:
    """
    (write (read . file "test.txt") . end "")
    """
    And a file named "test.txt" with:
    """
    foo
    bar
    baz
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    foo
    bar
    baz
    """
