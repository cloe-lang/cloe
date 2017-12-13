Feature: Read function
  Scenario: Read stdin
    Given a file named "main.coel" with:
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
    coel main.coel < test.txt
    """
    Then the stdout should contain exactly:
    """
    foo
    bar
    baz
    """

  Scenario: Read a file
    Given a file named "main.coel" with:
    """
    (write (read . file "test.txt") . end "")
    """
    And a file named "test.txt" with:
    """
    foo
    bar
    baz
    """
    When I successfully run `coel main.coel`
    Then the stdout should contain exactly:
    """
    foo
    bar
    baz
    """
