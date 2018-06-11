Feature: File System
  Scenario: Import file system module
    Given a file named "main.cloe" with:
    """
    (import "fs")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""

  Scenario: Create a directory
    Given a file named "main.cloe" with:
    """
    (import "fs")

    (fs.createDirectory "foo")
    """
    When I successfully run `cloe main.cloe`
    Then I successfully run `rmdir foo`

  Scenario: Remove an entry
    Given a file named "main.cloe" with:
    """
    (import "fs")

    (fs.remove "foo.txt")
    """
    And a file named "foo.txt" with ""
    When I successfully run `cloe main.cloe`
    Then I run `ls foo.txt`
    And the exit status should not be 0
